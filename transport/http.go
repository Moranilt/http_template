package transport

import (
	"net/http"
	"time"

	"github.com/Moranilt/http_template/endpoints"
	"github.com/Moranilt/http_template/middleware"
	"github.com/gorilla/mux"
)

func New(addr string, endpoints []endpoints.Endpoint, mw *middleware.Middleware) *http.Server {
	router := mux.NewRouter()
	router.Use(mw.Default)

	for _, endpoint := range endpoints {
		router.HandleFunc(endpoint.Pattern, endpoint.HandleFunc).Methods(endpoint.Methods...)
	}

	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  40 * time.Second,
		WriteTimeout: 40 * time.Second,
	}
	server.Handler = router
	return server
}
