package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Moranilt/http_template/clients/rabbitmq"
	"github.com/Moranilt/http_template/models"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const TracerName string = "repository"

type Repository struct {
	db       *sqlx.DB
	rabbitmq *rabbitmq.Client
	redis    *redis.Client
}

func New(db *sqlx.DB, rabbitmq *rabbitmq.Client, redis *redis.Client) *Repository {
	return &Repository{
		db:       db,
		rabbitmq: rabbitmq,
		redis:    redis,
	}
}

func (repo *Repository) Test(ctx context.Context, req *models.TestRequest) (*models.TestResponse, error) {
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

func (repo *Repository) Files(ctx context.Context, req *models.FileRequest) (*models.FileResponse, error) {
	if req == nil {
		return nil, errors.New("not valid request data")
	}
	var fileNames []string
	for _, f := range req.Files {
		fileNames = append(fileNames, f.Filename)
	}
	newCtx, span := otel.Tracer(TracerName).Start(ctx, "Files", trace.WithAttributes(
		attribute.String("Name", req.Name),
		attribute.StringSlice("FIles", fileNames),
		attribute.String("OneMoreFile", req.OneMoreFile.Filename),
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

	return &models.FileResponse{
		Name:        req.Name,
		Files:       req.Files,
		OneMoreFile: req.OneMoreFile,
	}, nil
}
