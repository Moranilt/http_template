package repository

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Moranilt/http_template/clients/database"
	database_mock "github.com/Moranilt/http_template/clients/database/mock"
	rabbitmq_mock "github.com/Moranilt/http_template/clients/rabbitmq/mock"
	redis_mock "github.com/Moranilt/http_template/clients/redis/mock"
	"github.com/Moranilt/http_template/logger"
	"github.com/Moranilt/http_template/models"
	"github.com/go-redis/redismock/v9"
)

type mockedRepository struct {
	repo         *Repository
	sqlMock      sqlmock.Sqlmock
	rabbitmqMock rabbitmq_mock.RabbitMQMocker
	redisMock    redismock.ClientMock
}

func mockRepository(t *testing.T) *mockedRepository {
	mockDb, sqlMock := database_mock.NewSQlMock(t)
	mockRabbitMQ := rabbitmq_mock.NewRabbitMQ(t)
	mockRedis, redisMock := redis_mock.New()
	mockLogger := logger.NewSlog(io.Discard)

	repo := New(&database.Client{mockDb}, mockRabbitMQ, mockRedis, mockLogger)

	return &mockedRepository{
		repo:         repo,
		sqlMock:      sqlMock,
		rabbitmqMock: mockRabbitMQ,
		redisMock:    redisMock,
	}
}

func makePointer(v string) *string {
	return &v
}

func TestCreateUser(t *testing.T) {
	// Mock repository dependencies
	mockedRepo := mockRepository(t)

	t.Run("Success", func(t *testing.T) {
		expectedUser := models.TestRequest{
			Firstname:  "John",
			Lastname:   "Doe",
			Patronymic: makePointer("Michael"),
		}
		expectedID := "1"
		exectedBody, _ := json.Marshal(expectedUser)
		// Set up mock expectations
		mockedRepo.sqlMock.ExpectQuery(regexp.QuoteMeta(QUERY_InsertUser)).
			WithArgs(expectedUser.Firstname, expectedUser.Lastname, expectedUser.Patronymic).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(expectedID))

		mockedRepo.redisMock.ExpectSet(expectedID, exectedBody, REDIS_TTL).SetVal(string(exectedBody))

		mockedRepo.rabbitmqMock.ExpectPush(exectedBody, nil)
		// Call Test()
		response, err := mockedRepo.repo.CreateUser(context.Background(), &expectedUser)
		if err != nil {
			t.Error(err)
		}

		// Assert expectations met
		if response.ID != expectedID {
			t.Errorf("Expected ID %s, got %s", expectedID, response.ID)
		}
	})

	t.Run("Error query", func(t *testing.T) {
		// Set error expectations
		expectedUser := models.TestRequest{
			Firstname:  "John",
			Lastname:   "Doe",
			Patronymic: makePointer("Michael"),
		}
		expectedID := "1"

		expectedError := errors.New("query error")
		mockedRepo.sqlMock.ExpectQuery(regexp.QuoteMeta(QUERY_InsertUser)).
			WithArgs(expectedUser.Firstname, expectedUser.Lastname, expectedUser.Patronymic).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(expectedID)).WillReturnError(expectedError)

		response, err := mockedRepo.repo.CreateUser(context.Background(), &expectedUser)
		if err.Error() != expectedError.Error() {
			t.Errorf("Expected error %q but got %q", expectedError, err)
		}

		if response != nil {
			t.Errorf("Expected nil response on error, got %v", response)
		}
	})

	t.Run("Error redis set", func(t *testing.T) {
		// Set error expectations
		expectedUser := models.TestRequest{
			Firstname:  "John",
			Lastname:   "Doe",
			Patronymic: makePointer("Michael"),
		}
		expectedID := "1"
		exectedBody, _ := json.Marshal(expectedUser)
		expectedError := errors.New("redis error")
		// Set up mock expectations
		mockedRepo.sqlMock.ExpectQuery(regexp.QuoteMeta(QUERY_InsertUser)).
			WithArgs(expectedUser.Firstname, expectedUser.Lastname, expectedUser.Patronymic).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(expectedID))

		mockedRepo.redisMock.ExpectSet(expectedID, exectedBody, REDIS_TTL).SetErr(expectedError)

		response, err := mockedRepo.repo.CreateUser(context.Background(), &expectedUser)
		if err.Error() != expectedError.Error() {
			t.Errorf("Expected error %q but got %q", expectedError, err)
		}

		if response != nil {
			t.Errorf("Expected nil response on error, got %v", response)
		}
	})

	t.Run("Error rabbitmq push", func(t *testing.T) {
		// Set error expectations
		expectedUser := models.TestRequest{
			Firstname:  "John",
			Lastname:   "Doe",
			Patronymic: makePointer("Michael"),
		}
		expectedID := "1"
		exectedBody, _ := json.Marshal(expectedUser)
		expectedError := errors.New("rabbitmq error")
		// Set up mock expectations
		mockedRepo.sqlMock.ExpectQuery(regexp.QuoteMeta(QUERY_InsertUser)).
			WithArgs(expectedUser.Firstname, expectedUser.Lastname, expectedUser.Patronymic).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(expectedID))

		mockedRepo.redisMock.ExpectSet(expectedID, exectedBody, REDIS_TTL).SetVal(string(exectedBody))
		mockedRepo.rabbitmqMock.ExpectPush(exectedBody, expectedError)

		response, err := mockedRepo.repo.CreateUser(context.Background(), &expectedUser)
		if err.Error() != expectedError.Error() {
			t.Errorf("Expected error %q but got %q", expectedError, err)
		}

		if response != nil {
			t.Errorf("Expected nil response on error, got %v", response)
		}
	})

}
