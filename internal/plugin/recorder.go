package plugin

import (
        "context"
        "time"

        "github.com/sun-yryr/Rec-adio/internal/db/models"
)

// PlatformInfo はプラットフォームの情報を表す構造体です
type PlatformInfo struct {
        // ID はプラットフォームの一意な識別子です
        ID string

        // Name はプラットフォームの表示名です
        Name string

        // Version はプラットフォームのバージョンです
        Version string

        // URL はプラットフォームのウェブサイトURLです
        URL string

        // Features はプラットフォームがサポートする機能のリストです
        Features []string

        // ScheduleSupport はスケジュール取得をサポートするかどうかを示します
        ScheduleSupport bool

        // LiveDetectionSupport はライブ配信検出をサポートするかどうかを示します
        LiveDetectionSupport bool

        // PerformerSearchSupport は演者検索をサポートするかどうかを示します
        PerformerSearchSupport bool

        // AccountSearchSupport はアカウント検索をサポートするかどうかを示します
        AccountSearchSupport bool
}

// RecordingStatus は録音の状態を表す列挙型です
type RecordingStatus string

const (
        RecordingStatusPending   RecordingStatus = "pending"   // 録音待ち
        RecordingStatusRecording RecordingStatus = "recording" // 録音中
        RecordingStatusCompleted RecordingStatus = "completed" // 録音完了
        RecordingStatusFailed    RecordingStatus = "failed"    // 録音失敗
        RecordingStatusCanceled  RecordingStatus = "canceled"  // 録音キャンセル
)

// RecordingSession は録音セッションを表すインターフェースです
type RecordingSession interface {
        // GetID は録音セッションのIDを返します
        GetID() string

        // GetStatus は録音の現在のステータスを返します
        GetStatus() RecordingStatus

        // Stop は録音を停止します
        Stop(ctx context.Context) error

        // GetOutputPath は録音ファイルのパスを返します（完了後）
        GetOutputPath() string

        // GetMetadata は録音のメタデータを返します
        GetMetadata() map[string]interface{}

        // OnStatusChange はステータス変更時のコールバックを設定します
        OnStatusChange(callback func(status RecordingStatus))
}

// ScheduleFilter はスケジュール検索のフィルタを表す構造体です
type ScheduleFilter struct {
        StartTimeFrom time.Time
        StartTimeTo   time.Time
        Platform      string
        Title         string
        PerformerID   string
        AccountID     string
}

// RecorderPlugin はレコーダープラグインのインターフェースです
type RecorderPlugin interface {
        // GetPlatformInfo はサポートするプラットフォームの情報を返します
        GetPlatformInfo() PlatformInfo

        // Record は指定されたスケジュールに基づいて録音を開始します
        Record(ctx context.Context, schedule *models.Schedule) (RecordingSession, error)

        // GetSchedule はプラットフォームから番組スケジュールを取得します
        GetSchedule(ctx context.Context, filter ScheduleFilter) ([]*models.Schedule, error)

        // SearchByPerformer は演者に基づいて番組を検索します
        SearchByPerformer(ctx context.Context, performer *models.Performer) ([]*models.Schedule, error)

        // SearchByPlatformAccount はプラットフォームアカウントに基づいて配信を検索します
        SearchByPlatformAccount(ctx context.Context, account *models.PlatformAccount) ([]*models.Schedule, error)

        // Initialize はプラグインの初期化を行います
        Initialize(ctx context.Context, config map[string]interface{}) error

        // Shutdown はプラグインのシャットダウン処理を行います
        Shutdown(ctx context.Context) error
}