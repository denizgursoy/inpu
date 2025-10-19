package inpu

import (
	"context"
	"fmt"
	"strings"
)

type stringBufferLogger struct {
	errorBuffer strings.Builder
	warnBuffer  strings.Builder
	infoBuffer  strings.Builder
	debugBuffer strings.Builder
}

func newStringBufferLogger() *stringBufferLogger {
	return &stringBufferLogger{
		errorBuffer: strings.Builder{},
		warnBuffer:  strings.Builder{},
		infoBuffer:  strings.Builder{},
		debugBuffer: strings.Builder{},
	}
}

func (s *stringBufferLogger) Error(ctx context.Context, err error, msg string, fields ...any) {
	s.errorBuffer.WriteString(fmt.Sprintf(msg, fields...))
}

func (s *stringBufferLogger) Warn(ctx context.Context, msg string, fields ...any) {
	s.warnBuffer.WriteString(fmt.Sprintf(msg, fields...))
}

func (s *stringBufferLogger) Info(ctx context.Context, msg string, fields ...any) {
	s.infoBuffer.WriteString(fmt.Sprintf(msg, fields...))
}

func (s *stringBufferLogger) Debug(ctx context.Context, msg string, fields ...any) {
	s.debugBuffer.WriteString(fmt.Sprintf(msg, fields...))
}
