// Package main は、Recotodのエントリーポイント
package main

import (
	"fmt"
	"net"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/sun-yryr/recoto/internal/adapter/grpcserver"
	"github.com/sun-yryr/recoto/internal/adapter/grpcserver/middleware"
	"github.com/sun-yryr/recoto/internal/broker"
	"github.com/sun-yryr/recoto/internal/broker/nats"
	"github.com/sun-yryr/recoto/internal/config"
	"github.com/sun-yryr/recoto/internal/logger"
	healthv1 "github.com/sun-yryr/recoto/pkg/api/recoto/health/v1"
)

func main() {
	// Load configuration
	cfg := lo.Must(config.Load())

	// Initialize appLogger
	appLogger := lo.Must(logger.NewLogger(cfg))
	defer func() {
		if err := appLogger.Sync(); err != nil {
			appLogger.Error("failed to sync logger", zap.Error(err))
		}
	}()

	appLogger.Info("Starting server...")

	// Initialize NATS broker
	embBroker, err := nats.NewEmbeddedBroker(appLogger)
	if err != nil {
		appLogger.Fatal("failed to create broker", zap.Error(err))
	}

	// Initialize server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		appLogger.Fatal("failed to listen", zap.Error(err))
	}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.NewLoggerInterceptor(appLogger)),
	)
	healthv1.RegisterHealthServiceServer(
		srv,
		grpcserver.NewHealthService(embBroker, broker.NewHealthCheck(embBroker)),
	)
	reflection.Register(srv)

	// Start server
	if err := srv.Serve(lis); err != nil {
		appLogger.Fatal("failed to serve", zap.Error(err))
	}
}
