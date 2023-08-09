package logger

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"log/slog"

	"github.com/sirupsen/logrus"
)

type ContextKey string

const (
	CtxRequestId ContextKey = "request_id"
)

type SLogger struct {
	*slog.Logger
}

func NewSlog(output io.Writer) *SLogger {
	l := slog.New(slog.NewJSONHandler(output, nil))
	logger := &SLogger{
		l,
	}
	return logger
}

func (l *SLogger) Fatal(msg string, args ...any) {
	l.Error(msg, args)
	os.Exit(1)
}

func (l *SLogger) Errorf(format string, args ...any) {
	l.Error(fmt.Sprintf(format, args...))
}

func (l *SLogger) Debugf(format string, args ...any) {
	l.Debug(fmt.Sprintf(format, args...))
}

func (l *SLogger) Infof(format string, args ...any) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l *SLogger) WithRequestInfo(r *http.Request) *SLogger {
	l = l.WithRequestId(r.Context())
	var clientIP string

	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		clientIP = ip
	}

	return &SLogger{
		l.With(
			"path", r.URL.Path,
			"method", r.Method,
			"ip", clientIP,
		),
	}
}
func (l *SLogger) WithRequestId(ctx context.Context) *SLogger {
	requestId := ctx.Value(CtxRequestId)
	if requestId != "" {
		return &SLogger{
			l.With("request_id", requestId),
		}
	}
	return l
}

type Logger struct {
	*logrus.Logger
}

func New() *Logger {
	logger := &Logger{}
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Level = logrus.TraceLevel
	logger.Out = os.Stdout

	return logger
}

func (l *Logger) WithRequestInfo(r *http.Request) *Logger {
	requestId := r.Context().Value(CtxRequestId)

	fields := logrus.Fields{
		"request_id": requestId,
		"path":       r.URL.Path,
		"method":     r.Method,
	}

	if clientIP, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		fields["ip"] = clientIP
	}

	return &Logger{
		l.WithFields(fields).Logger,
	}
}
