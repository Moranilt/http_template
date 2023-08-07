package db

import (
	"context"

	"github.com/Moranilt/http_template/clients/credentials"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func New(ctx context.Context, driverName string, creds *credentials.DB, production bool) (*sqlx.DB, error) {
	db, err := sqlx.Open(driverName, creds.SourceString(production))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
