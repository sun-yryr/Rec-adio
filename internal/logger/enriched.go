package logger

import (
	"context"

	"go.uber.org/zap"

	"github.com/sun-yryr/recoto/internal/trace"
)

// FromContextWithTrace は context からロガーを取得し、TraceIDで拡張する。
func FromContextWithTrace(ctx context.Context) *zap.Logger {
	logger := FromContext(ctx)

	if traceID, ok := trace.IDFromContext(ctx); ok {
		logger = logger.With(zap.String("trace_id", string(traceID)))
	}

	return logger
}

// WithTraceID はロガーにTraceIDフィールドを追加する。
func WithTraceID(logger *zap.Logger, traceID trace.ID) *zap.Logger {
	return logger.With(zap.String("trace_id", string(traceID)))
}
