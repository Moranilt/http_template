package endpoints

import (
	"net/http"

	"git.zonatelecom.ru/fsin/censor/service"
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
	}
}
