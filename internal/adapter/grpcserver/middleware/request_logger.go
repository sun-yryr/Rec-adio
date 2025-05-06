// Package middleware は、gRPCサーバーのミドルウェアを提供する
package middleware

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewRequestLoggerInterceptor は、gRPCリクエストをログに記録するUnaryServerInterceptorを返す。
func NewRequestLoggerInterceptor(zapLogger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now()

		// リクエスト情報をログに出力
		zapLogger.Info("gRPC request received",
			zap.String("method", info.FullMethod),
			zap.Any("request", req),
		)

		// ハンドラを実行
		resp, err := handler(ctx, req)

		// レスポンス情報をログに出力
		duration := time.Since(startTime)
		if err != nil {
			st, _ := status.FromError(err)
			zapLogger.Error("gRPC request failed",
				zap.String("method", info.FullMethod),
				zap.String("status_code", st.Code().String()),
				zap.String("status_message", st.Message()),
				zap.Duration("duration", duration),
				zap.Error(err),
			)
		} else {
			zapLogger.Info("gRPC request completed",
				zap.String("method", info.FullMethod),
				zap.String("status_code", codes.OK.String()),
				zap.Duration("duration", duration),
				zap.Any("response", resp),
			)
		}

		return resp, err
	}
}

// NewStreamRequestLoggerInterceptor は、ストリーミングリクエスト用のロガーインターセプターを返す。
func NewStreamRequestLoggerInterceptor(zapLogger *zap.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		startTime := time.Now()
		
		// ストリーミングリクエスト開始をログに出力
		zapLogger.Info("gRPC streaming request started",
			zap.String("method", info.FullMethod),
			zap.Bool("is_client_stream", info.IsClientStream),
			zap.Bool("is_server_stream", info.IsServerStream),
		)

		// ハンドラを実行
		err := handler(srv, stream)

		// ストリーミングリクエスト終了をログに出力
		duration := time.Since(startTime)
		if err != nil {
			st, _ := status.FromError(err)
			zapLogger.Error("gRPC streaming request failed",
				zap.String("method", info.FullMethod),
				zap.String("status_code", st.Code().String()),
				zap.String("status_message", st.Message()),
				zap.Duration("duration", duration),
				zap.Error(err),
			)
		} else {
			zapLogger.Info("gRPC streaming request completed",
				zap.String("method", info.FullMethod),
				zap.String("status_code", codes.OK.String()),
				zap.Duration("duration", duration),
			)
		}

		return err
	}
}