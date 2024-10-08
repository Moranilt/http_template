package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Moranilt/http-utils/clients/database"
	"github.com/Moranilt/http-utils/clients/rabbitmq"
	"github.com/Moranilt/http-utils/clients/redis"
	"github.com/Moranilt/http-utils/clients/vault"
	"github.com/Moranilt/http-utils/logger"
	"github.com/Moranilt/http_template/config"
	"github.com/Moranilt/http_template/endpoints"
	"github.com/Moranilt/http_template/middleware"
	"github.com/Moranilt/http_template/repository"
	"github.com/Moranilt/http_template/service"
	"github.com/Moranilt/http_template/tracer"
	"github.com/Moranilt/http_template/transport"
	_ "github.com/golang-migrate/migrate/source/file"
	"golang.org/x/sync/errgroup"
)

const (
	RABBITMQ_QUEUE_NAME = "test_queue"

	DB_DRIVER_NAME = "postgres"
)

func Run(ctx context.Context) {
	log := logger.New(os.Stdout, logger.TYPE_JSON)
	logger.SetDefault(log)

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// Vault
	err = vault.Init(&vault.Config{
		MountPath: cfg.Vault.MountPath,
		Token:     cfg.Vault.Token,
		Host:      cfg.Vault.Host,
	})
	if err != nil {
		log.Fatalf("vault: %v", err)
	}

	// Database
	dbCreds, err := vault.GetCreds[database.Credentials](ctx, cfg.Vault.DbCredsPath)
	if err != nil {
		log.Fatalf("get db creds from vault: %v", err)
	}

	db, err := database.New(ctx, DB_DRIVER_NAME, dbCreds)
	if err != nil {
		log.Fatalf("db connection: %v", err)
	}
	defer db.Close()

	// Tracer
	tp, err := tracer.NewProvider(cfg.Tracer.URL, cfg.Tracer.Name)
	if err != nil {
		log.Fatalf("tracer: %v", err)
	}

	// RabbitMQ
	rabbitMQCreds, err := vault.GetCreds[rabbitmq.Credentials](ctx, cfg.Vault.RabbitMQCreds)
	if err != nil {
		log.Fatalf("get rabbitmq creds from vault: %v", err)
	}

	rabbitmqClient := rabbitmq.Init(ctx, RABBITMQ_QUEUE_NAME, log, rabbitMQCreds)
	rabbitmq.ReadMsgs(ctx, 5, 5*time.Second, ConsumeMessage)

	// Redis
	redisCreds, err := vault.GetCreds[redis.Credentials](ctx, cfg.Vault.RedisCreds)
	if err != nil {
		log.Fatalf("get redis creds from vault: %v", err)
	}

	redisClient, err := redis.New(ctx, redisCreds)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}

	repo := repository.New(db, rabbitmqClient, redisClient, log)
	svc := service.New(log, repo)
	mw := middleware.New(log)
	ep := endpoints.MakeEndpoints(svc, mw)
	health := endpoints.MakeHealth(db, rabbitmqClient, redisClient)
	ep = append(ep, health)
	server := transport.New(fmt.Sprintf(":%s", cfg.Port), ep, mw)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		<-gCtx.Done()
		return tp.Shutdown(context.Background())
	})

	g.Go(func() error {
		<-gCtx.Done()
		return server.Shutdown(context.Background())
	})

	g.Go(func() error {
		<-gCtx.Done()
		return rabbitmq.Close()
	})

	g.Go(func() error {
		return server.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Infof("exit with: %s", err)
	}
}

// Example of consume
// You can provide any logic in this callback for each received message
//
// Requeue received message
// If it was requeued already - just Ack this
func ConsumeMessage(ctx context.Context, d rabbitmq.RabbitDelivery) error {
	if d.Redelivered() {
		d.Ack(true)
		fmt.Println("Message redelivered: ", string(d.Body()), d.DeliveryTag())
		return errors.New("message redelivered")
	}

	fmt.Println("Message: ", string(d.Body()), d.DeliveryTag(), d.MessageCount())
	d.Nack(false, true)
	return nil
}
