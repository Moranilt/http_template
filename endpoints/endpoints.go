package endpoints

import (
	"net/http"

	"github.com/Moranilt/http-utils/clients/database"
	"github.com/Moranilt/http-utils/clients/rabbitmq"
	"github.com/Moranilt/http-utils/clients/redis"
	"github.com/Moranilt/http_template/healthcheck"
	"github.com/Moranilt/http_template/middleware"
	"github.com/Moranilt/http_template/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Endpoint struct {
	Pattern    string
	HandleFunc http.HandlerFunc
	Methods    []string
	Middleware []middleware.EndpointMiddlewareFunc
}

func MakeEndpoints(service service.Service, mw *middleware.Middleware) []Endpoint {
	return []Endpoint{
		{
			Pattern:    "/user",
			HandleFunc: service.CreateUser,
			Methods:    []string{http.MethodPost},
		},
		{
			Pattern:    "/files",
			HandleFunc: service.Files,
			Methods:    []string{http.MethodPost},
			Middleware: []middleware.EndpointMiddlewareFunc{mw.AppTokenRequired},
		},
		{
			Pattern:    "/random-number",
			HandleFunc: service.GetRandomNumber,
			Methods:    []string{http.MethodGet},
		},
		{
			Pattern:    "/metrics",
			HandleFunc: promhttp.Handler().ServeHTTP,
			Methods:    []string{http.MethodGet},
		},
	}
}

func MakeHealth(db *database.Client, rabbitmq rabbitmq.RabbitMQClient, redis *redis.Client) Endpoint {
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
