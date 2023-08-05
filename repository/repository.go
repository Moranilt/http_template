package repository

import (
	"context"
	"encoding/json"

	"github.com/Moranilt/http_template/clients/rabbitmq"
	"github.com/Moranilt/http_template/models"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const TracerName string = "repository"

type Repository struct {
	db       *sqlx.DB
	rabbitmq *rabbitmq.Client
}

func New(db *sqlx.DB, rabbitmq *rabbitmq.Client) *Repository {
	return &Repository{
		db:       db,
		rabbitmq: rabbitmq,
	}
}

func (repo *Repository) Test(ctx context.Context, req *models.TestReq) (*models.TestResponse, error) {
	newCtx, span := otel.Tracer(TracerName).Start(ctx, "Test", trace.WithAttributes(
		attribute.String("Name", req.Name),
	))
	defer span.End()

	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	err = repo.rabbitmq.Push(newCtx, b)
	if err != nil {
		return nil, err
	}

	return &models.TestResponse{
		Name: req.Name,
	}, nil
}
