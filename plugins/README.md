# plugins - Rec-adio プラグイン実装

このディレクトリには、Rec-adioのプラグイン実装が含まれています。

## ディレクトリ構造

- **recorders/**: レコーダープラグイン
  - **agqr/**: 超A&G+レコーダー
  - **twitter_space/**: Twitter Spaceレコーダー
  - **bilibili/**: BiliBili Liveレコーダー
- **notifications/**: 通知プラグイン
  - **line/**: LINE通知プラグイン
  - **webhook/**: Webhook通知プラグイン

## プラグイン開発ガイド

### レコーダープラグインの開発

レコーダープラグインを開発するには、`pkg/plugin/recorder`パッケージで定義されている`RecorderPlugin`インターフェースを実装します。

```go
package myrecorder

import (
	"context"

	"github.com/username/rec-adio/pkg/plugin/recorder"
	"github.com/username/rec-adio/pkg/types"
)

// MyRecorder は新しいレコーダープラグインの実装です
type MyRecorder struct {
	config Config
}

// Config はプラグインの設定です
type Config struct {
	APIKey     string
	Timeout    int
	Quality    string
}

// New は新しいMyRecorderインスタンスを作成します
func New() *MyRecorder {
	return &MyRecorder{}
}

// GetPlatformInfo はサポートするプラットフォームの情報を返します
func (r *MyRecorder) GetPlatformInfo() types.PlatformInfo {
	return types.PlatformInfo{
		ID:                    "myplatform",
		Name:                  "My Platform",
		Version:               "1.0.0",
		URL:                   "https://example.com",
		Features:              []string{"recording", "schedule"},
		ScheduleSupport:       true,
		LiveDetectionSupport:  true,
		PerformerSearchSupport: true,
		AccountSearchSupport:  true,
	}
}

// Record は指定されたスケジュールに基づいて録音を開始します
func (r *MyRecorder) Record(ctx context.Context, schedule types.Schedule) (recorder.RecordingSession, error) {
	// 録音処理の実装
	// ...
}

// GetSchedule はプラットフォームから番組スケジュールを取得します
func (r *MyRecorder) GetSchedule(ctx context.Context, filter types.ScheduleFilter) ([]types.Schedule, error) {
	// スケジュール取得処理の実装
	// ...
}

// SearchByPerformer は演者に基づいて番組を検索します
func (r *MyRecorder) SearchByPerformer(ctx context.Context, performer types.Performer) ([]types.Schedule, error) {
	// 演者ベースの検索処理の実装
	// ...
}

// SearchByPlatformAccount はプラットフォームアカウントに基づいて配信を検索します
func (r *MyRecorder) SearchByPlatformAccount(ctx context.Context, account types.PlatformAccount) ([]types.Schedule, error) {
	// アカウントベースの検索処理の実装
	// ...
}

// Initialize はプラグインの初期化を行います
func (r *MyRecorder) Initialize(ctx context.Context, config recorder.Config) error {
	// 設定の取得と検証
	apiKey, ok := config.Get("api_key").(string)
	if !ok || apiKey == "" {
		return errors.New("api_key is required")
	}

	timeout, ok := config.Get("timeout").(int)
	if !ok {
		timeout = 30 // デフォルト値
	}

	quality, ok := config.Get("quality").(string)
	if !ok {
		quality = "high" // デフォルト値
	}

	r.config = Config{
		APIKey:  apiKey,
		Timeout: timeout,
		Quality: quality,
	}

	return nil
}

// Shutdown はプラグインのシャットダウン処理を行います
func (r *MyRecorder) Shutdown(ctx context.Context) error {
	// シャットダウン処理の実装
	// ...
	return nil
}
```

### 通知プラグインの開発

通知プラグインを開発するには、`pkg/plugin/notification`パッケージで定義されている`EventListener`インターフェースを実装します。

```go
package mynotification

import (
	"github.com/username/rec-adio/pkg/plugin/notification"
	"github.com/username/rec-adio/pkg/types"
)

// MyNotification は新しい通知プラグインの実装です
type MyNotification struct {
	config Config
	enabled bool
}

// Config はプラグインの設定です
type Config struct {
	APIKey string
	Events []string
}

// New は新しいMyNotificationインスタンスを作成します
func New() *MyNotification {
	return &MyNotification{
		enabled: true,
	}
}

// OnEvent はイベント発生時に呼び出されるメソッドです
func (n *MyNotification) OnEvent(event types.Event) error {
	if !n.enabled {
		return nil
	}

	// イベントタイプに基づいてメッセージを作成
	var message string
	switch event.Type {
	case "recording.started":
		title, _ := event.Data["title"].(string)
		message = fmt.Sprintf("録音開始: %s", title)
	case "recording.completed":
		title, _ := event.Data["title"].(string)
		message = fmt.Sprintf("録音完了: %s", title)
	case "recording.failed":
		title, _ := event.Data["title"].(string)
		errorMsg, _ := event.Data["error"].(string)
		message = fmt.Sprintf("録音失敗: %s\nエラー: %s", title, errorMsg)
	default:
		message = fmt.Sprintf("イベント: %s", event.Type)
	}

	// 通知の送信
	// ...

	return nil
}

// GetSupportedEvents はリスナーがサポートするイベントタイプのリストを返します
func (n *MyNotification) GetSupportedEvents() []string {
	return n.config.Events
}

// Initialize はプラグインの初期化を行います
func (n *MyNotification) Initialize(config notification.Config) error {
	// 設定の取得と検証
	apiKey, ok := config.Get("api_key").(string)
	if !ok || apiKey == "" {
		return errors.New("api_key is required")
	}

	events, ok := config.Get("events").([]string)
	if !ok {
		events = []string{"recording.completed", "recording.failed"} // デフォルト値
	}

	n.config = Config{
		APIKey: apiKey,
		Events: events,
	}

	return nil
}
