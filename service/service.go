package service

import (
	"net/http"

	"github.com/Moranilt/http_template/logger"
	"github.com/Moranilt/http_template/models"
	"github.com/Moranilt/http_template/repository"
	"github.com/Moranilt/http_template/utils/handler"
)

type Service interface {
	Test(http.ResponseWriter, *http.Request)
}

type service struct {
	log  *logger.Logger
	repo *repository.Repository
}

func New(log *logger.Logger, repo *repository.Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s *service) Test(w http.ResponseWriter, r *http.Request) {
	handler.New[*models.TestReq, *models.TestResponse](w, r, s.log, s.repo.Test).WithQuery().Run(http.StatusOK, http.StatusBadRequest)
}
