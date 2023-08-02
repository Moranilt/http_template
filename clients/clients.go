package clients

import (
	"context"

	"git.zonatelecom.ru/fsin/censor/clients/credentials"
	"github.com/jmoiron/sqlx"
)

func DB(ctx context.Context, production bool, creds *credentials.DBCreds) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", creds.SourceString(production))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
