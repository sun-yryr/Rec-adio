package main

import (
	"fmt"

	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/sun-yryr/recoto/api/server"
	"github.com/sun-yryr/recoto/api/server/router/health"
	"github.com/sun-yryr/recoto/internal/broker"
	"github.com/sun-yryr/recoto/internal/broker/nats"
	"github.com/sun-yryr/recoto/internal/config"
	"github.com/sun-yryr/recoto/internal/logger"
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
	server := server.NewServer(
		health.NewHealthRouter(embBroker, logger, broker.NewBrokerHealthCheck(embBroker)),
	)

	// Start server
	server.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
