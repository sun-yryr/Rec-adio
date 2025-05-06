package logger

import (
	"context"

	"go.uber.org/zap"
)

type ctxKeyLogger struct{}

// WithLogger は、コンテキストにロガーを追加する。
func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger{}, logger)
}

// FromContext は、コンテキストからロガーを取得する。
func FromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(ctxKeyLogger{}).(*zap.Logger)
	if !ok {
		return zap.NewNop()
	}

	return logger
}
