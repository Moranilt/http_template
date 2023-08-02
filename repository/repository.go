package repository

import (
	"context"

	"github.com/Moranilt/http_template/clients"
	"github.com/Moranilt/http_template/models"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const TracerName string = "repository"

type Repository struct {
	db     *sqlx.DB
	rabbit *clients.RabbitMQClient
}

func New(db *sqlx.DB, rabbit *clients.RabbitMQClient) *Repository {
	return &Repository{
		db:     db,
		rabbit: rabbit,
	}
}

func (repo *Repository) Test(ctx context.Context, req *models.TestReq) (*models.TestResponse, error) {
	newCtx, span := otel.Tracer(TracerName).Start(ctx, "Test", trace.WithAttributes(
		attribute.String("Name", req.Name),
	))
	defer span.End()

	err := repo.rabbit.Push(newCtx, []byte("Hello World!"))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
