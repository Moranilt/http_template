package service

import (
	"net/http"

	"github.com/Moranilt/http_template/logger"
	"github.com/Moranilt/http_template/repository"
	"github.com/Moranilt/http_template/utils/handler"
)

type Service interface {
	Test(http.ResponseWriter, *http.Request)
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

func (s *service) Test(w http.ResponseWriter, r *http.Request) {
	handler.New(w, r, s.log, s.repo.Test).
		WithQuery().
		Run(http.StatusOK, http.StatusBadRequest)
}

func (s *service) Files(w http.ResponseWriter, r *http.Request) {
	handler.New(w, r, s.log, s.repo.Files).
		WithMultipart(32<<20).
		Run(http.StatusOK, http.StatusBadRequest)
}
