package inpu

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type ctxKey string

const loggerCtxKey ctxKey = "inpu_logger"

type LogLevel int

const (
	LogLevelError LogLevel = iota
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

var DefaultLogger = NewLogger(LogLevelInfo)

func ContextWithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}

func GetLoggerFromContext(ctx context.Context) Logger {
	logger, ok := ctx.Value(loggerCtxKey).(Logger)
	if !ok {
		return DefaultLogger // fallback
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
	logger *slog.Logger
}

func NewLogger(level LogLevel) Logger {
	lvl := new(slog.LevelVar)
	lvl.Set(covertToSlogLevel(level))

	return &slogInpuLogger{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: lvl,
		})),
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

func covertToSlogLevel(lvl LogLevel) slog.Level {
	switch lvl {
	case LogLevelError:
		return slog.LevelError
	case LogLevelInfo:
		return slog.LevelInfo
	case LogLevelDebug:
		return slog.LevelDebug
	case LogLevelWarn:
		return slog.LevelWarn
	}

	return slog.LevelDebug
}
