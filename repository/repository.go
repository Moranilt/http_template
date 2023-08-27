package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Moranilt/http_template/clients/database"
	"github.com/Moranilt/http_template/clients/rabbitmq"
	"github.com/Moranilt/http_template/clients/redis"
	"github.com/Moranilt/http_template/logger"
	"github.com/Moranilt/http_template/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const TracerName string = "repository"

type Repository struct {
	db       *database.Client
	rabbitmq rabbitmq.RabbitMQClient
	redis    *redis.Client
	log      *logger.SLogger
}

func New(db *database.Client, rabbitmq rabbitmq.RabbitMQClient, redis *redis.Client, logger *logger.SLogger) *Repository {
	return &Repository{
		db:       db,
		rabbitmq: rabbitmq,
		redis:    redis,
		log:      logger,
	}
}

func (repo *Repository) Test(ctx context.Context, req *models.TestRequest) (*models.TestResponse, error) {
	repo.log.WithRequestId(ctx).InfoContext(ctx, TracerName, "data", req)
	newCtx, span := otel.Tracer(TracerName).Start(ctx, "Test", trace.WithAttributes(
		attribute.String("Name", req.Name),
	))
	defer span.End()

	err := repo.redis.Set(newCtx, "test", "name", 30*time.Second).Err()
	if err != nil {
		return nil, err
	}

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
	repo.log.WithRequestId(ctx).InfoContext(ctx, TracerName, "data", req)
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
