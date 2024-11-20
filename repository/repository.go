package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Moranilt/http-utils/clients/database"
	"github.com/Moranilt/http-utils/clients/rabbitmq"
	"github.com/Moranilt/http-utils/clients/redis"
	"github.com/Moranilt/http-utils/logger"
	"github.com/Moranilt/http-utils/tiny_errors"
	"github.com/Moranilt/http_template/custom_errors"
	"github.com/Moranilt/http_template/models"
	"github.com/Moranilt/http_template/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/rand"
)

const (
	QUERY_InsertUser = "INSERT INTO test (firstname, lastname, patronymic) VALUES ($1, $2, $3) RETURNING id"
)

const (
	REDIS_TTL = 30 * time.Second
)

const TracerName string = "repository"

type Repository struct {
	db       *database.Client
	rabbitmq rabbitmq.RabbitMQClient
	redis    *redis.Client
	log      logger.Logger
}

func New(db *database.Client, rabbitmq rabbitmq.RabbitMQClient, redis *redis.Client, logger logger.Logger) *Repository {
	return &Repository{
		db:       db,
		rabbitmq: rabbitmq,
		redis:    redis,
		log:      logger,
	}
}

func (repo *Repository) GetRandomNumber(ctx context.Context, req *models.GetRandomNumberRequest) (*models.GetRandomNumberResponse, tiny_errors.ErrorHandler) {
	repo.log.WithRequestId(ctx).InfoContext(ctx, TracerName, "data", req)
	_, span := otel.Tracer(TracerName).Start(ctx, "GetRandomNumber", trace.WithAttributes())
	defer span.End()

	errFields := utils.ValidateRequiredFields(
		utils.NewRequiredField("Min", req.Min),
		utils.NewRequiredField("Max", req.Max),
	)
	if errFields != nil {
		err := tiny_errors.New(custom_errors.ERR_CODE_REQUIRED_FIELD, errFields...)
		span.RecordError(err)
		span.SetStatus(codes.Error, "ValidateRequiredFields")
		return nil, err
	}

	randomNumber := rand.Intn(req.Max-req.Min) + req.Min
	return &models.GetRandomNumberResponse{
		Number: randomNumber,
	}, nil
}

func (repo *Repository) CreateUser(ctx context.Context, req *models.TestRequest) (*models.TestResponse, tiny_errors.ErrorHandler) {
	repo.log.WithRequestId(ctx).InfoContext(ctx, TracerName, "data", req)
	newCtx, span := otel.Tracer(TracerName).Start(ctx, "Test", trace.WithAttributes(
		attribute.String("Firstname", req.Firstname),
		attribute.String("Lastname", req.Lastname),
		attribute.String("Patronymic", *req.Patronymic),
	))
	defer span.End()

	row := repo.db.QueryRowxContext(newCtx, QUERY_InsertUser, req.Firstname, req.Lastname, *req.Patronymic)
	if row.Err() != nil {
		return nil, tiny_errors.New(custom_errors.ERR_CODE_Database, tiny_errors.Message(row.Err().Error()))
	}

	var lastInsertId string
	err := row.Scan(&lastInsertId)
	if err != nil {
		return nil, tiny_errors.New(custom_errors.ERR_CODE_Database, tiny_errors.Message(row.Err().Error()))
	}

	b, err := json.Marshal(req)
	if err != nil {
		return nil, tiny_errors.New(custom_errors.ERR_CODE_Marshal, tiny_errors.Message(err.Error()))
	}

	err = repo.redis.Set(newCtx, lastInsertId, b, REDIS_TTL).Err()
	if err != nil {
		return nil, tiny_errors.New(custom_errors.ERR_CODE_Redis, tiny_errors.Message(err.Error()))
	}

	err = repo.rabbitmq.Push(newCtx, b)
	if err != nil {
		return nil, tiny_errors.New(custom_errors.ERR_CODE_RabbitMQ, tiny_errors.Message(err.Error()))
	}

	return &models.TestResponse{
		ID: lastInsertId,
	}, nil
}

func (repo *Repository) Files(ctx context.Context, req *models.FileRequest) (*models.FileResponse, tiny_errors.ErrorHandler) {
	repo.log.WithRequestId(ctx).InfoContext(ctx, TracerName, "data", req)
	if req == nil {
		return nil, tiny_errors.New(custom_errors.ERR_CODE_BodyRequired)
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
		return nil, tiny_errors.New(custom_errors.ERR_CODE_Marshal, tiny_errors.Message(err.Error()))
	}

	err = repo.rabbitmq.Push(newCtx, b)
	if err != nil {
		return nil, tiny_errors.New(
			custom_errors.ERR_CODE_RabbitMQ,
			tiny_errors.Message(err.Error()),
			tiny_errors.Detail("event", "Push"),
		)
	}

	return &models.FileResponse{
		Name:        req.Name,
		Files:       req.Files,
		OneMoreFile: req.OneMoreFile,
	}, nil
}
