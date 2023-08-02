package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.zonatelecom.ru/fsin/censor/clients"
	"git.zonatelecom.ru/fsin/censor/clients/credentials"
	"git.zonatelecom.ru/fsin/censor/clients/vault"
	"git.zonatelecom.ru/fsin/censor/config"
	"git.zonatelecom.ru/fsin/censor/endpoints"
	"git.zonatelecom.ru/fsin/censor/logger"
	"git.zonatelecom.ru/fsin/censor/middleware"
	"git.zonatelecom.ru/fsin/censor/repository"
	"git.zonatelecom.ru/fsin/censor/service"
	"git.zonatelecom.ru/fsin/censor/tracer"
	"git.zonatelecom.ru/fsin/censor/transport"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"golang.org/x/sync/errgroup"
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

	db, err := clients.DB(ctx, cfg.Production, dbCreds)
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

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatal("migration: ", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", dbCreds.DBName, driver)
	if err != nil {
		log.Fatal("migration: ", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("migration: ", err)
	}

	version, _, err := m.Version()
	if err != nil {
		log.Fatal("migration: ", err)
	}

	log.Debug("migration: ", fmt.Sprintf("version %d", version))

	rabbitMQCreds, err := vault.GetCreds[credentials.RabbitMQCreds](ctx, cfg.Vault.RabbitMQCreds)
	if err != nil {
		log.Fatal("get rabbitmq creds from vault: ", err)
	}
	rabbitMQ, err := clients.RabbitMQ(ctx, log, rabbitMQCreds)
	if err != nil {
		log.Fatal("rabbitmq: ", err)
	}

	// msgCallback := func(d amqp.Delivery) error {
	// 	if d.Redelivered {
	// 		log.Println("Message redelivered...", d.DeliveryTag)
	// 		d.Ack(true)
	// 		return errors.New("Message redelivered...")
	// 	}
	// 	log.Println("Message: ", string(d.Body), d.DeliveryTag)
	// 	d.Nack(false, true)
	// 	return nil
	// }

	// _, err = rabbitMQ.Messages(ctx, 5, msgCallback)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// go func() {

	// 	i := 0
	// 	for d := range msgs {
	// 		if d.Redelivered {
	// 			log.Println("Message redelivered...", d.DeliveryTag)
	// 			d.Ack(true)
	// 			continue
	// 		}
	// 		log.Println("Message: ", string(d.Body), d.DeliveryTag)
	// 		d.Nack(false, true)

	// 		if i == 5 {
	// 			i = 0
	// 			<-time.After(15 * time.Second)
	// 		} else {
	// 			i++
	// 		}
	// 		// d.Nack(true, true)

	// 	}

	// 	// d.Ack(true)

	// }()

	repo := repository.New(db, rabbitMQ)
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
		return rabbitMQ.Close()
	})
	g.Go(func() error {
		return server.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Debugf("exit with: %s", err)
	}

}
