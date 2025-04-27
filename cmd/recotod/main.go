package main

import (
	"context"
	"log"

	"github.com/sun-yryr/recoto/internal/broker/nats"
	"github.com/sun-yryr/recoto/internal/config"
	"github.com/sun-yryr/recoto/internal/event/recording"
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

	ctx := context.Background()

	startedService := recording.NewStartedService(broker, logger)
	finishedService := recording.NewFinishedService(broker, logger)

	unsub, err := startedService.Subscribe(ctx, func(ctx context.Context, event *recording.StartedEvent) {
		logger.Info("Received started event", zap.String("recordingID", event.RecordingID), zap.Time("timestamp", event.Timestamp))
	})
	if err != nil {
		logger.Fatal("failed to subscribe to started events", zap.Error(err))
	}
	defer unsub()

	unsub2, err := finishedService.Subscribe(ctx, func(ctx context.Context, event *recording.FinishedEvent) {
		logger.Info("Received finished event", zap.String("recordingID", event.RecordingID), zap.Time("timestamp", event.Timestamp))
	})
	if err != nil {
		logger.Fatal("failed to subscribe to finished events", zap.Error(err))
	}
	defer unsub2()

	startedService.Publish(ctx, "sample uuid")

	// TODO: signalを処理する
	select {}
}
