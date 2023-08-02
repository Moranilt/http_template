package service

import (
	"net/http"

	"git.zonatelecom.ru/fsin/censor/logger"
	"git.zonatelecom.ru/fsin/censor/models"
	"git.zonatelecom.ru/fsin/censor/repository"
	"git.zonatelecom.ru/fsin/censor/utils/handler"
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
