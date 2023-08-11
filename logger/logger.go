package logger

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sync/atomic"

	"log/slog"

	"github.com/sirupsen/logrus"
)

type ContextKey string

var defaultLogger atomic.Value

func init() {
	defaultLogger.Store(NewSlog(os.Stdout))
}

func SetDefault(l *SLogger) {
	defaultLogger.Store(l)
}

func Default() *SLogger {
	return defaultLogger.Load().(*SLogger)
}

const (
	CtxRequestId ContextKey = "request_id"
)

const (
	LevelTrace  = slog.Level(-8)
	LevelNotice = slog.Level(2)
	LevelFatal  = slog.Level(12)
)

var LevelNames = map[slog.Leveler]string{
	LevelTrace:  "TRACE",
	LevelNotice: "NOTICE",
	LevelFatal:  "FATAL",
}

type SLogger struct {
	*slog.Logger
}

func NewSlog(output io.Writer) *SLogger {
	l := slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level:       LevelTrace,
		ReplaceAttr: renameLevel,
	}))

	logger := &SLogger{
		l,
	}
	return logger
}

func (l *SLogger) Trace(msg string, args ...any) {
	l.Log(context.Background(), LevelTrace, msg, args...)
}

func (l *SLogger) Notice(msg string, args ...any) {
	l.Log(context.Background(), LevelNotice, msg, args...)
}

func (l *SLogger) Fatal(msg string, args ...any) {
	l.Log(context.Background(), LevelFatal, msg, args...)
	os.Exit(1)
}

func (l *SLogger) Fatalf(format string, args ...any) {
	l.Log(context.Background(), LevelFatal, fmt.Sprintf(format, args...))
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

func renameLevel(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		level := a.Value.Any().(slog.Level)
		levelLabel, exists := LevelNames[level]
		if !exists {
			levelLabel = level.String()
		}

		a.Value = slog.StringValue(levelLabel)
	}

	return a
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
