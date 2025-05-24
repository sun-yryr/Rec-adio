// Package middleware は、gRPCサーバーのミドルウェアを提供する
package middleware

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/sun-yryr/recoto/internal/logger"
	"github.com/sun-yryr/recoto/internal/trace"
)

// NewLoggerInterceptor は TraceID 対応のロガーインターセプターを返す。
func NewLoggerInterceptor(zapLogger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// TraceID を生成または取得
		ctx, traceID := trace.GetOrGenerate(ctx)

		// TraceID でロガーを拡張
		enhancedLogger := logger.WithTraceID(zapLogger, traceID)
		ctx = logger.WithLogger(ctx, enhancedLogger)

		// リクエスト開始ログ
		enhancedLogger.Info("gRPC request started",
			zap.String("method", info.FullMethod),
			zap.String("trace_id", string(traceID)),
		)

		resp, err := handler(ctx, req)

		// リクエスト終了ログ
		if err != nil {
			enhancedLogger.Error("gRPC request failed",
				zap.String("method", info.FullMethod),
				zap.Error(err),
			)
		} else {
			enhancedLogger.Info("gRPC request completed",
				zap.String("method", info.FullMethod),
			)
		}

		return resp, err
	}
}
