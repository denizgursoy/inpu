package zap

import (
	"context"
	"fmt"

	"github.com/denizgursoy/inpu"
	"go.uber.org/zap"
)

func NewInpuLoggerFromZapLogger(logger *zap.Logger) inpu.Logger {
	return &inpuZapLogger{
		logger: logger.Sugar(),
	}
}

type inpuZapLogger struct {
	logger *zap.SugaredLogger
}

func (i *inpuZapLogger) Error(ctx context.Context, err error, msg string, fields ...any) {
	requestID := inpu.ExtractRequestIDFromContext(ctx)
	if requestID != nil {
		i.logger.Errorw(fmt.Sprintf(msg, fields...), append([]any{inpu.LoggerKeyRequestID, *requestID, "error", err})...)
	} else {
		i.logger.Errorw(fmt.Sprintf(msg, fields...), append([]any{"error", err})...)
	}
}

func (i *inpuZapLogger) Warn(ctx context.Context, msg string, fields ...any) {
	requestID := inpu.ExtractRequestIDFromContext(ctx)
	if requestID != nil {
		i.logger.Warnw(fmt.Sprintf(msg, fields...), append([]any{inpu.LoggerKeyRequestID, *requestID})...)
	} else {
		i.logger.Warnw(fmt.Sprintf(msg, fields...))
	}
}

func (i *inpuZapLogger) Info(ctx context.Context, msg string, fields ...any) {
	requestID := inpu.ExtractRequestIDFromContext(ctx)
	if requestID != nil {
		i.logger.Infow(fmt.Sprintf(msg, fields...), append([]any{inpu.LoggerKeyRequestID, *requestID})...)
	} else {
		i.logger.Infow(fmt.Sprintf(msg, fields...))
	}
}

func (i *inpuZapLogger) Debug(ctx context.Context, msg string, fields ...any) {
	requestID := inpu.ExtractRequestIDFromContext(ctx)
	if requestID != nil {
		i.logger.Debugw(fmt.Sprintf(msg, fields...), append([]any{inpu.LoggerKeyRequestID, *requestID})...)
	} else {
		i.logger.Debugw(fmt.Sprintf(msg, fields...))
	}
}
