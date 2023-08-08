package database

import (
	"context"

	"github.com/Moranilt/http_template/clients/credentials"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Client struct {
	*sqlx.DB
}

func (d *Client) Check(ctx context.Context) error {
	return d.PingContext(ctx)
}

func New(ctx context.Context, driverName string, creds *credentials.DB, production bool) (*Client, error) {
	connection, err := sqlx.Open(driverName, creds.SourceString(production))
	if err != nil {
		return nil, err
	}

	if err := connection.Ping(); err != nil {
		return nil, err
	}

	return &Client{
		connection,
	}, nil
}
