package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/sun-yryr/Rec-adio/internal/api"
	"github.com/sun-yryr/Rec-adio/internal/config"
	"github.com/sun-yryr/Rec-adio/internal/db/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	// バージョン情報
	version = "0.1.0"

	// 設定ファイルのパス
	configPath string

	// 設定
	cfg *config.Config

	// コンテキスト
	ctx    context.Context
	cancel context.CancelFunc

	// 待機グループ
	wg sync.WaitGroup
)

// rootCmd は、アプリケーションのルートコマンドです
var rootCmd = &cobra.Command{
	Use:     "recadiod",
	Short:   "Rec-adio Daemon - ラジオ録音管理デーモン",
	Version: version,
	Long: `Rec-adio Daemon は、インターネットラジオの録音を管理するためのデーモンサービスです。
A&G+、Radiko、音泉、響、Twitter Space、BiliBiliなどの複数のプラットフォームに対応しています。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 設定ファイルを読み込む
		var err error
		cfg, err = config.LoadConfig(configPath)
		if err != nil {
			// 設定ファイルが存在しない場合は、デフォルト設定を使用
			if os.IsNotExist(err) {
				log.Printf("設定ファイルが見つかりません: %s", configPath)
				log.Printf("デフォルト設定を使用します")
				cfg = config.DefaultConfig()
			} else {
				return fmt.Errorf("設定ファイルの読み込みに失敗しました: %w", err)
			}
		}

		// データディレクトリが存在しない場合は作成
		if _, err := os.Stat(cfg.App.DataDir); os.IsNotExist(err) {
			if err := os.MkdirAll(cfg.App.DataDir, 0755); err != nil {
				return fmt.Errorf("データディレクトリの作成に失敗しました: %w", err)
			}
		}

		// 一時ディレクトリが存在しない場合は作成
		if _, err := os.Stat(cfg.App.TempDir); os.IsNotExist(err) {
			if err := os.MkdirAll(cfg.App.TempDir, 0755); err != nil {
				return fmt.Errorf("一時ディレクトリの作成に失敗しました: %w", err)
			}
		}

		// 録音ファイル保存ディレクトリが存在しない場合は作成
		if _, err := os.Stat(cfg.App.RecordingsDir); os.IsNotExist(err) {
			if err := os.MkdirAll(cfg.App.RecordingsDir, 0755); err != nil {
				return fmt.Errorf("録音ファイル保存ディレクトリの作成に失敗しました: %w", err)
			}
		}

		db, err := gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("データベースへの接続に失敗しました: %w", err)
		}

		db.AutoMigrate(
			&models.Performer{},
			&models.PlatformAccount{},
			&models.ProgramInfo{},
			&models.Schedule{},
			&models.Record{},
			&models.SchedulePerformer{},
			&models.ProgramPerformer{},
			&models.ProgramPlatformAccount{},
			&models.PerformerSubscription{},
			&models.PlatformAccountSubscription{},
		)

		// コンテキストを作成
		ctx, cancel = context.WithCancel(context.Background())

		// シグナルハンドラを設定
		setupSignalHandler()

		// デーモンを起動
		if err := runDaemon(); err != nil {
			return fmt.Errorf("デーモンの実行に失敗しました: %w", err)
		}

		return nil
	},
}

// initConfigCmd は、設定ファイルを初期化するコマンドです
var initConfigCmd = &cobra.Command{
	Use:   "init-config",
	Short: "設定ファイルを初期化します",
	Long:  `デフォルト設定で設定ファイルを初期化します。既存の設定ファイルは上書きされます。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// デフォルト設定を取得
		defaultCfg := config.DefaultConfig()

		// 設定ファイルを保存
		if err := config.SaveConfig(defaultCfg, configPath); err != nil {
			return fmt.Errorf("設定ファイルの保存に失敗しました: %w", err)
		}

		fmt.Printf("設定ファイルを初期化しました: %s\n", configPath)
		return nil
	},
}

// versionCmd は、バージョン情報を表示するコマンドです
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "バージョン情報を表示します",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Rec-adio Daemon バージョン %s\n", version)
	},
}

// setupSignalHandler は、シグナルハンドラを設定します
func setupSignalHandler() {
	// シグナルチャネルを作成
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// シグナルを受信したらコンテキストをキャンセル
	go func() {
		sig := <-sigCh
		log.Printf("シグナルを受信しました: %s", sig)
		log.Println("シャットダウンを開始します...")
		cancel()
	}()
}

// runDaemon は、デーモンを実行します
func runDaemon() error {
	log.Println("Rec-adio Daemon を起動しています...")

	// APIサーバーを起動
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("APIサーバーを起動しています...")
		api.New(ctx)
	}()

	// スケジュール更新ワーカーを起動
	wg.Add(1)
	go runScheduleUpdateWorker()

	// 録音チェックワーカーを起動
	wg.Add(1)
	go runRecordingCheckWorker()

	// クリーンアップワーカーを起動
	wg.Add(1)
	go runCleanupWorker()

	// 全てのワーカーが終了するまで待機
	wg.Wait()

	log.Println("Rec-adio Daemon を終了しました")
	return nil
}

// runScheduleUpdateWorker は、スケジュール更新ワーカーを実行します
func runScheduleUpdateWorker() {
	defer wg.Done()

	log.Println("スケジュール更新ワーカーを起動しました")

	// スケジュール更新間隔（秒）
	interval := 300
	if cfg != nil {
		// TODO: デーモン設定からスケジュール更新間隔を取得
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// 初回実行
	updateSchedules()

	for {
		select {
		case <-ctx.Done():
			log.Println("スケジュール更新ワーカーを終了します")
			return
		case <-ticker.C:
			updateSchedules()
		}
	}
}

// updateSchedules は、録音スケジュールを更新します
func updateSchedules() {
	log.Println("録音スケジュールを更新しています...")
	// TODO: 録音スケジュールの更新処理を実装
}

// runRecordingCheckWorker は、録音チェックワーカーを実行します
func runRecordingCheckWorker() {
	defer wg.Done()

	log.Println("録音チェックワーカーを起動しました")

	// 録音チェック間隔（秒）
	interval := 60
	if cfg != nil {
		// TODO: デーモン設定から録音チェック間隔を取得
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// 初回実行
	checkRecordings()

	for {
		select {
		case <-ctx.Done():
			log.Println("録音チェックワーカーを終了します")
			return
		case <-ticker.C:
			checkRecordings()
		}
	}
}

// checkRecordings は、録音をチェックします
func checkRecordings() {
	log.Println("録音をチェックしています...")
	// TODO: 録音チェック処理を実装
}

// runCleanupWorker は、クリーンアップワーカーを実行します
func runCleanupWorker() {
	defer wg.Done()

	log.Println("クリーンアップワーカーを起動しました")

	// クリーンアップ間隔（秒）
	interval := 86400 // 1日
	if cfg != nil {
		// TODO: デーモン設定からクリーンアップ間隔を取得
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// 初回実行
	cleanup()

	for {
		select {
		case <-ctx.Done():
			log.Println("クリーンアップワーカーを終了します")
			return
		case <-ticker.C:
			cleanup()
		}
	}
}

// cleanup は、古い録音ファイルなどをクリーンアップします
func cleanup() {
	log.Println("クリーンアップを実行しています...")
	// TODO: クリーンアップ処理を実装
}

func init() {
	// 設定ファイルのデフォルトパスを設定
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	defaultConfigPath := filepath.Join(homeDir, ".config", "recadio", "daemon.toml")

	// フラグを設定
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", defaultConfigPath, "設定ファイルのパス")

	// サブコマンドを追加
	rootCmd.AddCommand(initConfigCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
