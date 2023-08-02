package logger

import (
	"net"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type ContextKey string

const (
	CtxRequestId ContextKey = "request_id"
)

type Logger struct {
	logrus.Logger
}

func New() *Logger {
	logger := Logger{}
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Level = logrus.TraceLevel
	logger.Out = os.Stdout

	return &logger
}

func (l *Logger) WithRequestInfo(r *http.Request) *logrus.Entry {
	logger := New()
	requestId := r.Context().Value(CtxRequestId)

	fields := logrus.Fields{
		"request_id": requestId,
		"path":       r.URL.Path,
		"method":     r.Method,
	}

	if clientIP, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		fields["ip"] = clientIP
	}

	return logger.WithFields(fields)
}
