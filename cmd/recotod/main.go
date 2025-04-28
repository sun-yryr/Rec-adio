package main

import (
	"fmt"
	"log"

	"github.com/sun-yryr/recoto/api/server"
	"github.com/sun-yryr/recoto/api/server/router/health"
	"github.com/sun-yryr/recoto/internal/broker/nats"
	"github.com/sun-yryr/recoto/internal/config"
	"github.com/sun-yryr/recoto/internal/logger"
	"go.uber.org/zap"
)

func main() {
	// 設定の読み込み
	// TODO: .envの取り扱い
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// ロガーの設定
	logger, err := logger.NewLogger(cfg)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting server...")

	broker, err := nats.NewEmbeddedBroker()
	if err != nil {
		logger.Fatal("failed to create broker", zap.Error(err))
	}

	server := server.NewServer(
		health.NewHealthRouter(broker, logger),
	)

	server.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
