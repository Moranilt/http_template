package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Moranilt/http_template/clients/credentials"
	"github.com/Moranilt/http_template/clients/db"
	"github.com/Moranilt/http_template/clients/rabbitmq"
	"github.com/Moranilt/http_template/clients/vault"
	"github.com/Moranilt/http_template/config"
	"github.com/Moranilt/http_template/endpoints"
	"github.com/Moranilt/http_template/logger"
	"github.com/Moranilt/http_template/middleware"
	"github.com/Moranilt/http_template/repository"
	"github.com/Moranilt/http_template/service"
	"github.com/Moranilt/http_template/tracer"
	"github.com/Moranilt/http_template/transport"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
)

const (
	RABBITMQ_QUEUE_NAME = "test_queue"
)

func main() {
	log := logger.New()
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		<-exit
		cancel()
	}()

	cfg, err := config.Read()
	if err != nil {
		log.Fatal("config: ", err)
	}

	err = vault.Init(&vault.Config{
		MountPath: cfg.Vault.MountPath,
		Token:     cfg.Vault.Token,
		Host:      cfg.Vault.Host,
	})
	if err != nil {
		log.Fatal("vault: ", err)
	}

	dbCreds, err := vault.GetCreds[credentials.DBCreds](ctx, cfg.Vault.DbCredsPath)
	if err != nil {
		log.Fatal("get db creds from vault: ", err)
	}

	db, err := db.New(ctx, cfg.Production, dbCreds)
	if err != nil {
		log.Fatal("db connection: ", err)
	}
	defer db.Close()

	tp, err := tracer.NewProvider(cfg.Tracer.URL, cfg.Tracer.Name)
	if err != nil {
		log.Fatal("tracer: ", err)
	}
	defer func(ctx context.Context) {
		if err := tp.Shutdown(ctx); err != nil {
			log.Errorf("Error shutting down tracer provider: %v", err)
		}
	}(ctx)

	err = RunMigrations(log, db.DB, dbCreds.DBName)
	if err != nil {
		log.Fatal("migration: ", err)
	}

	// RabbitMQ
	rabbitMQCreds, err := vault.GetCreds[credentials.RabbitMQCreds](ctx, cfg.Vault.RabbitMQCreds)
	if err != nil {
		log.Fatal("get rabbitmq creds from vault: ", err)
	}

	rebbitmqClient := rabbitmq.Init(ctx, RABBITMQ_QUEUE_NAME, log, rabbitMQCreds)
	rabbitmq.ReadMsgs(ctx, 5, ConsumeMessage)

	repo := repository.New(db, rebbitmqClient)
	svc := service.New(log, repo)
	ep := endpoints.MakeEndpoints(svc)
	mw := middleware.New(log)
	server := transport.New(fmt.Sprintf(":%s", cfg.Port), ep, mw)

	g, gCtx := errgroup.WithContext(ctx)
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
		log.Debugf("exit with: %s", err)
	}

}

// Example of consume
// You can provide any logic in this callback for each received message
//
// Requeue received message
// If it was requeued already - just Ack this
func ConsumeMessage(d amqp.Delivery) error {
	if d.Redelivered {
		d.Ack(true)
		fmt.Println("Message redelivered: ", string(d.Body), d.DeliveryTag)
		return errors.New("message redelivered")
	}

	fmt.Println("Message: ", string(d.Body), d.DeliveryTag, d.MessageCount)
	d.Nack(false, true)
	return nil
}

func RunMigrations(log *logger.Logger, db *sql.DB, databaseName string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", databaseName, driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	version, _, err := m.Version()
	if err != nil {
		return err
	}

	log.Debug("migration: ", fmt.Sprintf("version %d", version))
	return nil
}
