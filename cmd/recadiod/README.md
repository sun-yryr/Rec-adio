# recadiod - Rec-adio コアサービス

このディレクトリには、Rec-adioのコアサービス（デーモン）の実装が含まれています。

## 概要

recadiodは、Rec-adioのバックグラウンドサービスとして動作し、以下の機能を提供します：

- APIサーバーの提供
- スケジュール管理
- 録音ジョブの実行
- プラグインの管理
- イベントの発行と処理

## 実装方法

コアサービスは、以下のコンポーネントで構成されています：

1. **APIサーバー**: RESTful APIを提供し、CLIツールからの要求を処理します
2. **スケジュールマネージャー**: 録音スケジュールを管理します
3. **録音マネージャー**: 実際の録音処理を担当します
4. **プラグインマネージャー**: レコーダープラグインと通知プラグインを管理します
5. **イベントシステム**: システム内のイベントを処理します
6. **設定マネージャー**: 設定を管理します

## サンプルコード

```go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/username/rec-adio/internal/api"
	"github.com/username/rec-adio/internal/config"
	"github.com/username/rec-adio/internal/core/event"
	"github.com/username/rec-adio/internal/core/recording"
	"github.com/username/rec-adio/internal/core/schedule"
	"github.com/username/rec-adio/internal/db"
)

func main() {
	// 設定の読み込み
	cfg, err := config.Load("config.toml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// データベースの初期化
	database, err := db.New(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// イベントシステムの初期化
	eventSystem := event.NewSystem()

	// マネージャーの初期化
	scheduleManager := schedule.NewManager(database, eventSystem)
	recordingManager := recording.NewManager(database, eventSystem)

	// APIサーバーの初期化
	apiServer := api.NewServer(cfg.API, scheduleManager, recordingManager, eventSystem)

	// シグナルハンドリング
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down...")
		cancel()
	}()

	// サーバー起動
	if err := apiServer.Start(ctx); err != nil {
		log.Fatalf("Failed to start API server: %v", err)
	}
}
