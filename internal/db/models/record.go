package models

import (
	"time"

	"gorm.io/gorm"
)

// RecordStatus は、録音のステータスを表す列挙型です
type RecordStatus string

const (
	RecordStatusPending   RecordStatus = "pending"   // 録音待ち
	RecordStatusRecording RecordStatus = "recording" // 録音中
	RecordStatusCompleted RecordStatus = "completed" // 録音完了
	RecordStatusFailed    RecordStatus = "failed"    // 録音失敗
	RecordStatusCanceled  RecordStatus = "canceled"  // 録音キャンセル
)

// Record は、録音結果の情報を表します
type Record struct {
	ID           string         `gorm:"primaryKey" json:"id"`              // 一意な識別子
	ScheduleID   string         `gorm:"not null;index" json:"schedule_id"` // 関連するスケジュールのID
	Title        string         `gorm:"not null;index" json:"title"`       // タイトル
	RecordTime   time.Time      `gorm:"not null;index" json:"record_time"` // 録音時間
	Duration     int            `json:"duration"`                          // 期間（秒）
	FilePath     string         `json:"file_path"`                         // ファイルパス
	FileSize     int64          `json:"file_size"`                         // ファイルサイズ（バイト）
	Platform     string         `gorm:"not null;index" json:"platform"`    // プラットフォーム
	Status       RecordStatus   `gorm:"not null;index" json:"status"`      // ステータス
	ErrorMessage string         `json:"error_message"`                     // エラーメッセージ（失敗時）
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`  // 作成日時
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`  // 更新日時
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // 削除日時（ソフトデリート用）

	// リレーション
	Schedule *Schedule `gorm:"foreignKey:ScheduleID" json:"-"`
}

// BeforeCreate は、レコード作成前に呼び出されるフック関数です
func (r *Record) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == "" {
		r.ID = GenerateRecordID()
	}
	return nil
}

// NewRecord は、新しいRecordインスタンスを作成します
func NewRecord(scheduleID, title string, recordTime time.Time, platform string) *Record {
	return &Record{
		ScheduleID: scheduleID,
		Title:      title,
		RecordTime: recordTime,
		Platform:   platform,
		Status:     RecordStatusPending,
	}
}

// SetDuration は、期間を設定します
func (r *Record) SetDuration(duration int) *Record {
	r.Duration = duration
	return r
}

// SetFilePath は、ファイルパスを設定します
func (r *Record) SetFilePath(filePath string) *Record {
	r.FilePath = filePath
	return r
}

// SetFileSize は、ファイルサイズを設定します
func (r *Record) SetFileSize(fileSize int64) *Record {
	r.FileSize = fileSize
	return r
}

// SetStatus は、ステータスを設定します
func (r *Record) SetStatus(status RecordStatus) *Record {
	r.Status = status
	return r
}

// SetErrorMessage は、エラーメッセージを設定します
func (r *Record) SetErrorMessage(errorMessage string) *Record {
	r.ErrorMessage = errorMessage
	return r
}

// IsCompleted は、録音が完了しているかどうかを返します
func (r *Record) IsCompleted() bool {
	return r.Status == RecordStatusCompleted
}

// IsFailed は、録音が失敗しているかどうかを返します
func (r *Record) IsFailed() bool {
	return r.Status == RecordStatusFailed
}

// IsCanceled は、録音がキャンセルされているかどうかを返します
func (r *Record) IsCanceled() bool {
	return r.Status == RecordStatusCanceled
}

// IsInProgress は、録音が進行中かどうかを返します
func (r *Record) IsInProgress() bool {
	return r.Status == RecordStatusPending || r.Status == RecordStatusRecording
}
