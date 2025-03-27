package db

import (
        "time"

        "github.com/sun-yryr/Rec-adio/internal/db/models"
)

// RecordFilter は録音レコードの検索フィルタを表します
type RecordFilter struct {
        ScheduleID    string
        Platform      string
        Status        models.RecordStatus
        StartTimeFrom time.Time
        StartTimeTo   time.Time
        Title         string
}

// ScheduleFilter はスケジュールの検索フィルタを表します
type ScheduleFilter struct {
        Platform      string
        StartTimeFrom time.Time
        StartTimeTo   time.Time
        Title         string
        PerformerID   string
        IsProcessing  *bool
}

// PerformerFilter は演者の検索フィルタを表します
type PerformerFilter struct {
        Name string
}

// PlatformAccountFilter はプラットフォームアカウントの検索フィルタを表します
type PlatformAccountFilter struct {
        Platform string
        Username string
}

// SubscriptionFilter は購読の検索フィルタを表します
type SubscriptionFilter struct {
        PerformerID string
        AccountID   string
        Platform    string
}