package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Moranilt/http-utils/logger"
	"github.com/Moranilt/http-utils/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type ContextKey string

const (
	TOKEN_HEADER = "X-App-Token"
)

type Middleware struct {
	logger               logger.Logger
	requestStatusCounter *prometheus.CounterVec
	responseTime         *prometheus.HistogramVec
}

type EndpointMiddlewareFunc func(handleFunc http.Handler) http.Handler

func New(l logger.Logger) *Middleware {
	return &Middleware{
		logger: l,
		requestStatusCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "your_app_http_response_total",
			Help: "Total number of requests by endpoint and status",
		},
			[]string{"method", "endpoint", "status"},
		),
		responseTime: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name: "your_app_http_response_time_seconds",
			Help: "Response time in seconds",
		},
			[]string{"method", "endpoint", "status"},
		),
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

func (m *Middleware) Otel(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// extract trace id from header
		newTraceID := r.Header.Get("X-Trace-ID")

		if newTraceID != "" {
			traceID, err := trace.TraceIDFromHex(newTraceID)
			if err == nil {
				spanContext := trace.NewSpanContext(trace.SpanContextConfig{
					TraceID: trace.TraceID(traceID),
				})
				// create new span context with trace id from header
				ctx = trace.ContextWithSpanContext(r.Context(), spanContext)
			}
		}

		path, _ := mux.CurrentRoute(r).GetPathTemplate()
		ctx, span := otel.Tracer("http").Start(ctx, path)
		defer span.End()

		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.route", path),
			attribute.String("span.kind", "server"),
			attribute.String("request_id", GetRequestID(r.Context())),
		)

		rw := &responseWriter{ResponseWriter: w}
		r = r.WithContext(ctx)

		traceID := span.SpanContext().TraceID().String()
		w.Header().Add("X-Trace-ID", traceID)
		otel.GetTextMapPropagator().Inject(r.Context(), propagation.HeaderCarrier(w.Header()))

		next.ServeHTTP(rw, r)

		span.SetAttributes(attribute.Int("http.status_code", rw.statusCode))
	})
}

func (m *Middleware) Prometheus(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w}
		r = r.WithContext(r.Context())
		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		path, _ := mux.CurrentRoute(r).GetPathTemplate()
		status := strconv.Itoa(rw.statusCode)
		m.requestStatusCounter.WithLabelValues(r.Method, path, status).Inc()
		m.responseTime.WithLabelValues(r.Method, path, status).Observe(duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (m *Middleware) AppTokenRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(TOKEN_HEADER) == "" {
			response.ErrorResponse(w, fmt.Errorf("%s required", TOKEN_HEADER), http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetRequestID(ctx context.Context) string {
	return ctx.Value(logger.CtxRequestId).(string)
}
