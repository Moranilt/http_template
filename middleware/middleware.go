package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Moranilt/http_template/logger"
	"github.com/google/uuid"
)

type ContextKey string

const (
	TOKEN_HEADER = "X-App-Token"

	CtxAppToken ContextKey = "app_token"
)

type Middleware struct {
	logger *logger.Logger
}

func GetAppTokenFromContext(ctx context.Context) (string, error) {
	ctxValue := ctx.Value(CtxAppToken)
	token, ok := ctxValue.(string)
	if !ok {
		return "", fmt.Errorf("not valid type of token. Got %T, wants %q", token, "string")
	}
	return token, nil
}

func New(l *logger.Logger) *Middleware {
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
