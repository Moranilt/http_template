package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Moranilt/http_template/logger"
	"github.com/Moranilt/http_template/utils/response"
	"github.com/google/uuid"
)

type ContextKey string

const (
	TOKEN_HEADER = "X-App-Token"
)

type Middleware struct {
	logger *logger.SLogger
}

type EndpointMiddlewareFunc func(handleFunc http.HandlerFunc) http.HandlerFunc

func New(l *logger.SLogger) *Middleware {
	return &Middleware{
		logger: l,
	}
}

func (m *Middleware) Default(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := uuid.NewString()
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("X-Request-ID", requestId)
		ctx := context.WithValue(r.Context(), logger.CtxRequestId, requestId)
		r = r.WithContext(ctx)
		m.logger.WithRequestInfo(r).Info("Incoming request")
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) AppTokenRequired(handleFunc http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(TOKEN_HEADER) == "" {
			response.ErrorResponse(w, fmt.Errorf("%s required", TOKEN_HEADER), http.StatusBadRequest)
			return
		}
		handleFunc(w, r)
	})
}
