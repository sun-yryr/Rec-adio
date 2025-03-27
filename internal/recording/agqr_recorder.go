package recording

import (
        "context"
        "fmt"
        "path/filepath"
        "time"

        "github.com/sun-yryr/Rec-adio/internal/db/models"
        "github.com/sun-yryr/Rec-adio/internal/plugin"
)

// AGQRRecorder は超A&G+のレコーダープラグインです
type AGQRRecorder struct {
        outputDir string
        streamURL string
}

// NewAGQRRecorder は新しいAGQRRecorderインスタンスを作成します
func NewAGQRRecorder(outputDir, streamURL string) *AGQRRecorder {
        return &AGQRRecorder{
                outputDir: outputDir,
                streamURL: streamURL,
        }
}

// GetPlatformInfo はサポートするプラットフォームの情報を返します
func (r *AGQRRecorder) GetPlatformInfo() plugin.PlatformInfo {
        return plugin.PlatformInfo{
                ID:                   "agqr",
                Name:                 "超A&G+",
                Version:              "1.0.0",
                URL:                  "https://www.agqr.jp/",
                Features:             []string{"recording", "schedule"},
                ScheduleSupport:      true,
                LiveDetectionSupport: false,
                PerformerSearchSupport: true,
                AccountSearchSupport:   false,
        }
}

// Record は指定されたスケジュールに基づいて録音を開始します
func (r *AGQRRecorder) Record(ctx context.Context, schedule *models.Schedule) (plugin.RecordingSession, error) {
        // 出力ディレクトリの作成
        outputDir := filepath.Join(r.outputDir, "agqr")

        // セッションの作成
        session := NewFFmpegSession(
                schedule.ID,
                r.streamURL,
                outputDir,
                schedule.Title,
                schedule.Duration,
        )

        // 録音の開始
        if err := session.Start(ctx); err != nil {
                return nil, fmt.Errorf("failed to start recording: %w", err)
        }

        return session, nil
}

// GetSchedule はプラットフォームから番組スケジュールを取得します
func (r *AGQRRecorder) GetSchedule(ctx context.Context, filter plugin.ScheduleFilter) ([]*models.Schedule, error) {
        // 実際の実装では、超A&G+のウェブサイトからスクレイピングなどでスケジュールを取得する
        // ここではダミーデータを返す
        schedules := []*models.Schedule{
                models.NewSchedule(
                        "サンプル番組1",
                        time.Now().Add(1*time.Hour),
                        3600, // 1時間
                        "agqr",
                ),
                models.NewSchedule(
                        "サンプル番組2",
                        time.Now().Add(3*time.Hour),
                        1800, // 30分
                        "agqr",
                ),
        }

        return schedules, nil
}

// SearchByPerformer は演者に基づいて番組を検索します
func (r *AGQRRecorder) SearchByPerformer(ctx context.Context, performer *models.Performer) ([]*models.Schedule, error) {
        // 実際の実装では、演者名に基づいて番組を検索する
        // ここではダミーデータを返す
        schedules := []*models.Schedule{
                models.NewSchedule(
                        fmt.Sprintf("%sのラジオ", performer.Name),
                        time.Now().Add(2*time.Hour),
                        3600, // 1時間
                        "agqr",
                ),
        }

        return schedules, nil
}

// SearchByPlatformAccount はプラットフォームアカウントに基づいて配信を検索します
func (r *AGQRRecorder) SearchByPlatformAccount(ctx context.Context, account *models.PlatformAccount) ([]*models.Schedule, error) {
        // 超A&G+ではアカウントベースの検索はサポートしていないため、空の結果を返す
        return []*models.Schedule{}, nil
}

// Initialize はプラグインの初期化を行います
func (r *AGQRRecorder) Initialize(ctx context.Context, config map[string]interface{}) error {
        // 設定からストリームURLを取得
        if url, ok := config["stream_url"].(string); ok && url != "" {
                r.streamURL = url
        }

        // 設定から出力ディレクトリを取得
        if dir, ok := config["output_dir"].(string); ok && dir != "" {
                r.outputDir = dir
        }

        return nil
}

// Shutdown はプラグインのシャットダウン処理を行います
func (r *AGQRRecorder) Shutdown(ctx context.Context) error {
        // 特に何もする必要がない
        return nil
}