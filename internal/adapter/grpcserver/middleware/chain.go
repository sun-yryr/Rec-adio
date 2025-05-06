// Package middleware は、gRPCサーバーのミドルウェアを提供する
package middleware

import (
        "context"

        "google.golang.org/grpc"
)

// ChainUnaryInterceptors は、複数のUnaryServerInterceptorを1つのインターセプターにチェーンする。
func ChainUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
        return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
                // 最後のインターセプターから順に実行するチェーンを構築
                chainHandler := handler
                for i := len(interceptors) - 1; i >= 0; i-- {
                        i := i // ループ変数をキャプチャ
                        chainHandler = func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
                                return interceptors[i](currentCtx, currentReq, info, chainHandler)
                        }
                }
                return chainHandler(ctx, req)
        }
}

// ChainStreamInterceptors は、複数のStreamServerInterceptorを1つのインターセプターにチェーンする。
func ChainStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
        return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
                // 最後のインターセプターから順に実行するチェーンを構築
                chainHandler := handler
                for i := len(interceptors) - 1; i >= 0; i-- {
                        i := i // ループ変数をキャプチャ
                        chainHandler = func(currentSrv interface{}, currentStream grpc.ServerStream) error {
                                return interceptors[i](currentSrv, currentStream, info, chainHandler)
                        }
                }
                return chainHandler(srv, stream)
        }
}