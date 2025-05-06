// Package middleware は、gRPCサーバーのミドルウェアを提供する
package middleware

import (
        "context"
        "time"

        "go.uber.org/zap"
        "google.golang.org/grpc"
        "google.golang.org/grpc/codes"
        "google.golang.org/grpc/status"

        "github.com/sun-yryr/recoto/internal/logger"
)

// NewLoggerInterceptor は、ctxにLoggerを追加するUnaryServerInterceptorを返す。
func NewLoggerInterceptor(zapLogger *zap.Logger) grpc.UnaryServerInterceptor {
        return func(
                ctx context.Context,
                req interface{},
                info *grpc.UnaryServerInfo,
                handler grpc.UnaryHandler,
        ) (interface{}, error) {
                ctx = logger.WithLogger(ctx, zapLogger)
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

// NewStreamLoggerInterceptor は、ストリーミングリクエスト用のロガーインターセプターを返す。
func NewStreamLoggerInterceptor(zapLogger *zap.Logger) grpc.StreamServerInterceptor {
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

                // コンテキストにロガーを追加したラッパーストリームを作成
                ctx := logger.WithLogger(stream.Context(), zapLogger)
                wrappedStream := &wrappedServerStream{
                        ServerStream: stream,
                        ctx:          ctx,
                }

                // ハンドラを実行
                err := handler(srv, wrappedStream)

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

// wrappedServerStream は、コンテキストを上書きするServerStreamラッパー
type wrappedServerStream struct {
        grpc.ServerStream
        ctx context.Context
}

// Context は、ラップされたコンテキストを返す
func (w *wrappedServerStream) Context() context.Context {
        return w.ctx
}
