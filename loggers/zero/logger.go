package zero

import (
	"context"

	"github.com/denizgursoy/inpu"
	"github.com/rs/zerolog/log"
)

func NewInpuLoggerFromZeroLog() inpu.Logger {
	return &inpuZeroLogger{}
}

type inpuZeroLogger struct{}

func (i *inpuZeroLogger) Error(ctx context.Context, err error, msg string, fields ...any) {
	requestID := inpu.ExtractRequestIDFromContext(ctx)
	if requestID != nil {
		log.Ctx(ctx).Error().Err(err).Str(inpu.LoggerKeyRequestID, *requestID).Msgf(msg, fields...)
	} else {
		log.Ctx(ctx).Error().Err(err).Msgf(msg, fields...)
	}
}

func (i *inpuZeroLogger) Warn(ctx context.Context, msg string, fields ...any) {
	requestID := inpu.ExtractRequestIDFromContext(ctx)
	if requestID != nil {
		log.Ctx(ctx).Warn().Str(inpu.LoggerKeyRequestID, *requestID).Msgf(msg, fields...)
	} else {
		log.Ctx(ctx).Warn().Msgf(msg, fields...)
	}
}

func (i *inpuZeroLogger) Info(ctx context.Context, msg string, fields ...any) {
	requestID := inpu.ExtractRequestIDFromContext(ctx)
	if requestID != nil {
		log.Ctx(ctx).Info().Str(inpu.LoggerKeyRequestID, *requestID).Msgf(msg, fields...)
	} else {
		log.Ctx(ctx).Info().Msgf(msg, fields...)
	}
}

func (i *inpuZeroLogger) Debug(ctx context.Context, msg string, fields ...any) {
	requestID := inpu.ExtractRequestIDFromContext(ctx)
	if requestID != nil {
		log.Ctx(ctx).Debug().Str(inpu.LoggerKeyRequestID, *requestID).Msgf(msg, fields...)
	} else {
		log.Ctx(ctx).Debug().Msgf(msg, fields...)
	}
}
