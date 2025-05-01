package main

import (
	"fmt"
	"net"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	mygrpc "github.com/sun-yryr/recoto/internal/adapter/grpc"
	"github.com/sun-yryr/recoto/internal/broker"
	"github.com/sun-yryr/recoto/internal/broker/nats"
	"github.com/sun-yryr/recoto/internal/config"
	"github.com/sun-yryr/recoto/internal/logger"
	healthv1 "github.com/sun-yryr/recoto/pkg/api/health/v1"
)

func main() {
	// Load configuration
	cfg := lo.Must(config.Load())

	// Initialize logger
	logger := lo.Must(logger.NewLogger(cfg))
	defer logger.Sync()

	logger.Info("Starting server...")

	// Initialize NATS broker
	embBroker, err := nats.NewEmbeddedBroker()
	if err != nil {
		logger.Fatal("failed to create broker", zap.Error(err))
	}

	// Initialize server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	healthv1.RegisterHealthServiceServer(
		grpcServer,
		mygrpc.NewHealthService(embBroker, logger, broker.NewBrokerHealthCheck(embBroker)),
	)
	reflection.Register(grpcServer)

	// Start server
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
