// Package main は、Recotodのエントリーポイント
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/sun-yryr/recoto/internal/adapter/grpcserver"
	"github.com/sun-yryr/recoto/internal/adapter/grpcserver/middleware"
	"github.com/sun-yryr/recoto/internal/broker"
	"github.com/sun-yryr/recoto/internal/broker/nats"
	config "github.com/sun-yryr/recoto/internal/config/server"
	"github.com/sun-yryr/recoto/internal/event/recording"
	"github.com/sun-yryr/recoto/internal/fileutil"
	"github.com/sun-yryr/recoto/internal/logger"
	"github.com/sun-yryr/recoto/internal/recorder"
	healthv1 "github.com/sun-yryr/recoto/pkg/api/recoto/health/v1"
	recordingv1 "github.com/sun-yryr/recoto/pkg/api/recoto/recording/v1"
)

//nolint:funlen // main関数は許して
func main() {
	// コマンドライン引数の解析
	var configPath string

	flag.StringVar(
		&configPath,
		"config",
		"",
		"Path to config file (default: $HOME/.config/recoto/daemon.toml)",
	)
	flag.Parse()

	// 設定を読み込む
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 設定のバリデーション
	if err := config.ValidateConfig(cfg); err != nil {
		log.Fatalf("Invalid config: %v", err)
	}

	// Initialize appLogger
	appLogger, err := logger.NewLogger(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	defer func() {
		if err := appLogger.Sync(); err != nil {
			appLogger.Error("failed to sync logger", zap.Error(err))
		}
	}()

	// saveDirを作成
	if err := os.MkdirAll(cfg.Recording.SaveDir, fileutil.DefaultDirPerm); err != nil {
		appLogger.Fatal("failed to create saveDir", zap.Error(err))
	}

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
	recordingRequestedService := recording.NewRequestedService(embBroker)
	recordingStartedService := recording.NewStartedService(embBroker)
	recordingFinishedService := recording.NewFinishedService(embBroker)

	// Initialize recorder
	urlRecorder := recorder.NewURLRecorder(appLogger)
	recordingManager := recorder.NewRecordingManager(
		recordingRequestedService,
		recordingStartedService,
		recordingFinishedService,
		appLogger,
	)

	// URLRecorderを追加
	if err := recordingManager.AddRecorder(urlRecorder); err != nil {
		// エラーログを出力するだけで、サーバー自体は起動する
		appLogger.Warn("failed to add URL recorder", zap.Error(err))
	}

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

	appLogger.Info("Server is listening on port", zap.Int("port", cfg.Server.Port))

	if err := srv.Serve(lis); err != nil {
		appLogger.Fatal("failed to serve", zap.Error(err))
	}
}
