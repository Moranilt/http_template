package endpoints

import (
	"net/http"

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
