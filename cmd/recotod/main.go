// Package main は、Recotodのエントリーポイント
package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/sun-yryr/recoto/internal/adapter/grpcserver"
	"github.com/sun-yryr/recoto/internal/adapter/grpcserver/middleware"
	"github.com/sun-yryr/recoto/internal/broker"
	"github.com/sun-yryr/recoto/internal/broker/nats"
	"github.com/sun-yryr/recoto/internal/config"
	"github.com/sun-yryr/recoto/internal/event/recording"
	"github.com/sun-yryr/recoto/internal/logger"
	"github.com/sun-yryr/recoto/internal/recorder"
	healthv1 "github.com/sun-yryr/recoto/pkg/api/recoto/health/v1"
	recordingv1 "github.com/sun-yryr/recoto/pkg/api/recoto/recording/v1"
)

const saveDirPerm = 0o750

//nolint:funlen // main関数は許して
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

	// saveDirを作成
	if err := os.MkdirAll(cfg.Recording.SaveDir, saveDirPerm); err != nil {
		appLogger.Fatal("failed to create saveDir", zap.Error(err))
	}

	appLogger.Info("Starting server...")

	// Initialize NATS broker
	embBroker, err := nats.NewEmbeddedBroker(appLogger)
	if err != nil {
		appLogger.Fatal("failed to create broker", zap.Error(err))
	}

	defer func() {
		if err := embBroker.Close(); err != nil {
			appLogger.Error("failed to close broker", zap.Error(err))
		}
	}()

	// イベント
	recordingRequestedService := recording.NewRequestedService(embBroker, appLogger)
	recordingStartedService := recording.NewStartedService(embBroker, appLogger)
	recordingFinishedService := recording.NewFinishedService(embBroker, appLogger)

	// Initialize recorder
	urlRecorder := recorder.NewURLRecorder(appLogger)
	recordingManager := recorder.NewRecordingManager(
		recordingRequestedService,
		recordingStartedService,
		recordingFinishedService,
		appLogger,
		[]recorder.Recorder{urlRecorder},
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := recordingManager.Start(ctx); err != nil {
		appLogger.Fatal("failed to start recording manager", zap.Error(err))
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
	recordingv1.RegisterRecordingServiceServer(
		srv,
		grpcserver.NewRecordingService(cfg, recordingRequestedService),
	)
	reflection.Register(srv)

	// Start server
	if err := srv.Serve(lis); err != nil {
		appLogger.Fatal("failed to serve", zap.Error(err))
	}
}
