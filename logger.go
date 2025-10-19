package inpu

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type ctxKey string

const loggerCtxKey ctxKey = "inpu_logger"

var defaultLogger = NewLogger(LogLevelSimple)

func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}

func FromContext(ctx context.Context) Logger {
	logger, ok := ctx.Value(loggerCtxKey).(Logger)
	if !ok {
		return defaultLogger // fallback
	}
	return logger
}

type Logger interface {
	Error(ctx context.Context, err error, msg string, fields ...any)
	Warn(ctx context.Context, msg string, fields ...any)
	Info(ctx context.Context, msg string, fields ...any)
	Debug(ctx context.Context, msg string, fields ...any)
}

type slogInpuLogger struct {
	level  LogLevel
	logger *slog.Logger
}

func NewLogger(level LogLevel) Logger {
	return &slogInpuLogger{
		level:  level,
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func (d *slogInpuLogger) Error(ctx context.Context, err error, msg string, fields ...any) {
	d.logger.ErrorContext(ctx, fmt.Sprintf(msg, fields...), "error", err)
}

func (d *slogInpuLogger) Warn(ctx context.Context, msg string, fields ...any) {
	d.logger.WarnContext(ctx, fmt.Sprintf(msg, fields...))
}

func (d *slogInpuLogger) Info(ctx context.Context, msg string, fields ...any) {
	d.logger.InfoContext(ctx, fmt.Sprintf(msg, fields...))
}

func (d *slogInpuLogger) Debug(ctx context.Context, msg string, fields ...any) {
	d.logger.DebugContext(ctx, fmt.Sprintf(msg, fields...))
}
