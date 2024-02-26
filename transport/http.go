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
		handlerFunc := applyMiddleware(endpoint.HandleFunc, endpoint.Middleware)
		router.HandleFunc(endpoint.Pattern, handlerFunc).Methods(endpoint.Methods...)
	}

	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  40 * time.Second,
		WriteTimeout: 40 * time.Second,
		Handler:      router,
	}
	return server
}

func applyMiddleware(handlerFunc http.HandlerFunc, mws []middleware.EndpointMiddlewareFunc) http.HandlerFunc {
	for _, mw := range mws {
		handlerFunc = mw(handlerFunc)
	}
	return handlerFunc
}
