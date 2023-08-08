package endpoints

import (
	"net/http"

	db "github.com/Moranilt/http_template/clients/database"
	"github.com/Moranilt/http_template/clients/rabbitmq"
	"github.com/Moranilt/http_template/clients/redis"
	"github.com/Moranilt/http_template/healthcheck"
	"github.com/Moranilt/http_template/service"
)

type Endpoint struct {
	Pattern    string
	HandleFunc http.HandlerFunc
	Methods    []string
}

func MakeEndpoints(service service.Service) []Endpoint {
	return []Endpoint{
		{
			Pattern:    "/test",
			HandleFunc: service.Test,
			Methods:    []string{http.MethodGet},
		},
		{
			Pattern:    "/files",
			HandleFunc: service.Files,
			Methods:    []string{http.MethodPost},
		},
	}
}

func MakeHealth(db *db.Client, rabbitmq *rabbitmq.Client, redis *redis.Client) Endpoint {
	return Endpoint{
		Pattern: "/health",
		HandleFunc: healthcheck.HandlerFunc(
			healthcheck.HealthItem{
				Name:    "database",
				Checker: db,
			},
			healthcheck.HealthItem{
				Name:    "rabbitmq",
				Checker: rabbitmq,
			},
			healthcheck.HealthItem{
				Name:    "redis",
				Checker: redis,
			},
		),
		Methods: []string{http.MethodGet},
	}
}
