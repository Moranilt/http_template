package service

import (
	"net/http"

	"github.com/Moranilt/http_template/logger"
	"github.com/Moranilt/http_template/repository"
	"github.com/Moranilt/http_template/utils/handler"
)

type Service interface {
	CreateUser(http.ResponseWriter, *http.Request)
	Files(w http.ResponseWriter, r *http.Request)
}

type service struct {
	log  *logger.SLogger
	repo *repository.Repository
}

func New(log *logger.SLogger, repo *repository.Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s *service) CreateUser(w http.ResponseWriter, r *http.Request) {
	handler.New(w, r, s.log, s.repo.CreateUser).
		WithJson().
		Run(http.StatusOK)
}

func (s *service) Files(w http.ResponseWriter, r *http.Request) {
	handler.New(w, r, s.log, s.repo.Files).
		WithMultipart(32 << 20).
		Run(http.StatusOK)
}
