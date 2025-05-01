// Package middleware は、gRPCサーバーのミドルウェアを提供します.
package middleware

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/sun-yryr/recoto/internal/logger"
)

// NewLoggerInterceptor は、ctxにLoggerを追加するUnaryServerInterceptorを返す.
func NewLoggerInterceptor(zapLogger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = logger.WithLogger(ctx, zapLogger)

		return handler(ctx, req)
	}
}
