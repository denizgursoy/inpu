package inpu

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

const LoggerKeyRequestID = "request_id"

type LogLevel int

const (
	LogLevelError LogLevel = iota
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

// DefaultLogger is the fallback logger used when no logger is found in the context.
// By default it is a no-op logger that discards all output. Assign SlogLogger or a
// custom Logger implementation to enable logging globally.
var DefaultLogger Logger = noopLogger{}

// SlogLogger is a ready-to-use JSON logger that writes to stdout at INFO level.
// Assign it to DefaultLogger or inject it via ContextWithLogger to enable logging.
var SlogLogger = NewInpuLoggerFromSlog(LogLevelInfo)

// noopLogger discards all log messages.
type noopLogger struct{}

func (noopLogger) Error(context.Context, error, string, ...any) {}
func (noopLogger) Warn(context.Context, string, ...any)         {}
func (noopLogger) Info(context.Context, string, ...any)         {}
func (noopLogger) Debug(context.Context, string, ...any)        {}

func ContextWithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, ContextKeyLogger, logger)
}

func ExtractLoggerFromContext(ctx context.Context) Logger {
	logger, ok := ctx.Value(ContextKeyLogger).(Logger)
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

func NewInpuLoggerFromSlog(level LogLevel) Logger {
	lvl := new(slog.LevelVar)
	lvl.Set(covertToSlogLevel(level))

	return &slogInpuLogger{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: lvl,
		})),
	}
}

func (d *slogInpuLogger) Error(ctx context.Context, err error, msg string, fields ...any) {
	requestID := ExtractRequestIDFromContext(ctx)
	if requestID != nil {
		d.logger.ErrorContext(ctx, fmt.Sprintf(msg, fields...), "error", err, LoggerKeyRequestID, *requestID)
	} else {
		d.logger.ErrorContext(ctx, fmt.Sprintf(msg, fields...), "error", err)
	}
}

func (d *slogInpuLogger) Warn(ctx context.Context, msg string, fields ...any) {
	requestID := ExtractRequestIDFromContext(ctx)
	if requestID != nil {
		d.logger.WarnContext(ctx, fmt.Sprintf(msg, fields...), LoggerKeyRequestID, *requestID)
	} else {
		d.logger.WarnContext(ctx, fmt.Sprintf(msg, fields...))
	}
}

func (d *slogInpuLogger) Info(ctx context.Context, msg string, fields ...any) {
	requestID := ExtractRequestIDFromContext(ctx)
	if requestID != nil {
		d.logger.InfoContext(ctx, fmt.Sprintf(msg, fields...), LoggerKeyRequestID, *requestID)
	} else {
		d.logger.InfoContext(ctx, fmt.Sprintf(msg, fields...))
	}
}

func (d *slogInpuLogger) Debug(ctx context.Context, msg string, fields ...any) {
	requestID := ExtractRequestIDFromContext(ctx)
	if requestID != nil {
		d.logger.DebugContext(ctx, fmt.Sprintf(msg, fields...), LoggerKeyRequestID, *requestID)
	} else {
		d.logger.DebugContext(ctx, fmt.Sprintf(msg, fields...))
	}
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
