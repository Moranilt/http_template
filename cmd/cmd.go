package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Moranilt/http-utils/clients/database"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

type CommandType string

const (
	UP   = CommandType("up")
	DOWN = CommandType("down")
	RUN  = CommandType("run")

	DB_DRIVER_NAME = "postgres"
)

type CliData struct {
	database.Credentials
	version string
	command CommandType
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		<-exit
		cancel()
	}()

	// Get database credentials from CLI
	cliData := getCliData()

	db, err := database.New(ctx, DB_DRIVER_NAME, &cliData.Credentials)
	if err != nil {
		log.Fatalf("db connection: %v", err)
	}
	defer db.Close()

	// Run database migrations
	if err := runMigrations(db.DB.DB, cliData.DBName, cliData.command, cliData.version); err != nil {
		log.Fatalf("database migrations failed: %v", err)
	}
}

func getCliData() *CliData {

	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("user", "root", "Database user")
	dbPass := flag.String("pass", "", "Database password")
	dbHost := flag.String("host", "localhost", "Database host")
	sslMode := flag.String("sslmode", "", "Database sslmode")
	version := flag.String("version", "latest", "Migration version")

	flag.Usage = func() {
		fmt.Print("\nCommands: \n")
		fmt.Print("  up - run database migrations up for one step\n")
		fmt.Print("  down - run database migrations down for one step\n")
		fmt.Print("  run - run database migrations up until the selected version\n")
		fmt.Println("-----------------------------------------------------------------------")
		flag.PrintDefaults()
	}
	arg := os.Args[1]

	switch arg {
	case string(UP):
	case string(DOWN):
	case string(RUN):
		break
	default:
		flag.Usage()
		log.Fatal("Command must be one of 'up', 'down', or 'run'")
	}

	os.Args = os.Args[1:]
	flag.Parse()

	if dbName == nil || *dbName == "" {
		flag.Usage()
		log.Fatal("Database name must be provided")
	}

	if dbPass == nil || *dbPass == "" {
		flag.Usage()
		log.Fatal("Database password must be provided")
	}

	if dbUser == nil || *dbUser == "" {
		flag.Usage()
		log.Fatal("Database user must be provided")
	}

	if dbHost == nil || *dbHost == "" {
		flag.Usage()
		log.Fatal("Database host must be provided")
	}

	if _, err := strconv.Atoi(*version); *version != "latest" && err != nil {
		flag.Usage()
		log.Fatal("Migration version must be a number or 'latest'")
	}

	dbCreds := database.Credentials{
		Username: *dbUser,
		Password: *dbPass,
		DBName:   *dbName,
		Host:     *dbHost,
		SSLMode:  sslMode,
	}

	cliData := &CliData{
		Credentials: dbCreds,
		version:     *version,
		command:     CommandType(arg),
	}

	return cliData
}

func runMigrations(db *sql.DB, databaseName string, command CommandType, newVersion string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", databaseName, driver)
	if err != nil {
		return err
	}

	switch command {
	case UP:
		err = m.Steps(1)
		if err != nil {
			return err
		}
	case DOWN:
		err = m.Steps(-1)
		if err != nil {
			return err
		}
	case RUN:
		if newVersion == "latest" {
			err = m.Up()
			if err != nil && err != migrate.ErrNoChange {
				return err
			}
		} else {
			intVersion, err := strconv.Atoi(newVersion)
			if err != nil {
				return err
			}
			err = m.Migrate(uint(intVersion))
			if err != nil && err != migrate.ErrNoChange {
				return err
			}
		}
	}

	version, _, err := m.Version()
	if err != nil {
		return err
	}

	log.Printf("migration: version %d\n", version)
	return nil
}
