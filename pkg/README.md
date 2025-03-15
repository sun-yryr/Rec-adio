# pkg - Rec-adio 公開パッケージ

このディレクトリには、Rec-adioの公開パッケージが含まれています。これらのパッケージは、Rec-adioの外部からの利用も想定されています。

## ディレクトリ構造

- **client/**: APIクライアント
- **plugin/**: プラグインAPI
  - **recorder/**: レコーダープラグインインターフェース
  - **notification/**: 通知プラグインインターフェース
- **types/**: 共通型定義

## プラグインAPI

Rec-adioは、プラグインシステムを通じて拡張可能です。プラグインAPIは、以下のインターフェースを提供します：

### レコーダープラグイン

レコーダープラグインは、特定のプラットフォームからの録音を担当します。

```go
// RecorderPlugin はレコーダープラグインのインターフェースです
type RecorderPlugin interface {
    // GetPlatformInfo はサポートするプラットフォームの情報を返します
    GetPlatformInfo() PlatformInfo
    
    // Record は指定されたスケジュールに基づいて録音を開始します
    Record(ctx context.Context, schedule Schedule) (RecordingSession, error)
    
    // GetSchedule はプラットフォームから番組スケジュールを取得します
    GetSchedule(ctx context.Context, filter ScheduleFilter) ([]Schedule, error)
    
    // SearchByPerformer は演者に基づいて番組を検索します
    SearchByPerformer(ctx context.Context, performer Performer) ([]Schedule, error)
    
    // SearchByPlatformAccount はプラットフォームアカウントに基づいて配信を検索します
    SearchByPlatformAccount(ctx context.Context, account PlatformAccount) ([]Schedule, error)
    
    // Initialize はプラグインの初期化を行います
    Initialize(ctx context.Context, config Config) error
    
    // Shutdown はプラグインのシャットダウン処理を行います
    Shutdown(ctx context.Context) error
}
```

### 通知プラグイン

通知プラグインは、システムイベントを外部に通知します。

```go
// EventListener はイベントリスナーのインターフェースです
type EventListener interface {
    // OnEvent はイベント発生時に呼び出されるメソッドです
    OnEvent(event Event) error
    
    // GetSupportedEvents はリスナーがサポートするイベントタイプのリストを返します
    GetSupportedEvents() []string
}
