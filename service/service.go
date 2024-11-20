package service

import (
	"net/http"

	"github.com/Moranilt/http-utils/handler"
	"github.com/Moranilt/http-utils/logger"
	"github.com/Moranilt/http_template/repository"
)

type Service interface {
	CreateUser(http.ResponseWriter, *http.Request)
	Files(w http.ResponseWriter, r *http.Request)
	GetRandomNumber(w http.ResponseWriter, r *http.Request)
}

type service struct {
	log  logger.Logger
	repo *repository.Repository
}

func New(log logger.Logger, repo *repository.Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s *service) CreateUser(w http.ResponseWriter, r *http.Request) {
	handler.New(w, r, s.log, s.repo.CreateUser).
		WithJSON().
		Run(http.StatusOK)
}

func (s *service) Files(w http.ResponseWriter, r *http.Request) {
	handler.New(w, r, s.log, s.repo.Files).
		WithMultipart(32 << 20).
		Run(http.StatusOK)
}

func (s *service) GetRandomNumber(w http.ResponseWriter, r *http.Request) {
	handler.New(w, r, s.log, s.repo.GetRandomNumber).
		WithQuery().
		Run(http.StatusOK)
}
