package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Schedule は、録音予定の情報を表します
type Schedule struct {
	ID            string          `gorm:"primaryKey" json:"id"`                              // 一意な識別子
	Title         string          `gorm:"not null;index" json:"title"`                       // タイトル
	StartTime     time.Time       `gorm:"not null;index" json:"start_time"`                  // 開始時間
	Duration      int             `gorm:"not null" json:"duration"`                          // 期間（秒）
	Platform      string          `gorm:"not null;index" json:"platform"`                    // プラットフォーム
	ProgramInfoID string          `gorm:"index" json:"program_info_id"`                      // 関連する番組情報のID（オプション）
	IsProcessing  bool            `gorm:"not null;default:false;index" json:"is_processing"` // 処理中かどうか
	ExtraData     json.RawMessage `json:"extra_data"`                                        // 追加データ（JSON）
	CreatedAt     time.Time       `gorm:"autoCreateTime" json:"created_at"`                  // 作成日時
	UpdatedAt     time.Time       `gorm:"autoUpdateTime" json:"updated_at"`                  // 更新日時
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`                 // 削除日時（ソフトデリート用）

	// リレーション
	ProgramInfo        *ProgramInfo         `gorm:"foreignKey:ProgramInfoID" json:"-"`
	Records            []*Record            `gorm:"foreignKey:ScheduleID" json:"-"`
	Performers         []*Performer         `gorm:"many2many:schedule_performers;" json:"-"`
	SchedulePerformers []*SchedulePerformer `gorm:"foreignKey:ScheduleID" json:"-"`
}

// BeforeCreate は、レコード作成前に呼び出されるフック関数です
func (s *Schedule) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = GenerateScheduleID()
	}
	return nil
}

// NewSchedule は、新しいScheduleインスタンスを作成します
func NewSchedule(title string, startTime time.Time, duration int, platform string) *Schedule {
	return &Schedule{
		Title:        title,
		StartTime:    startTime,
		Duration:     duration,
		Platform:     platform,
		IsProcessing: false,
	}
}

// SetProgramInfoID は、関連する番組情報のIDを設定します
func (s *Schedule) SetProgramInfoID(programInfoID string) *Schedule {
	s.ProgramInfoID = programInfoID
	return s
}

// SetIsProcessing は、処理中かどうかを設定します
func (s *Schedule) SetIsProcessing(isProcessing bool) *Schedule {
	s.IsProcessing = isProcessing
	return s
}

// SetExtraData は、追加データを設定します
func (s *Schedule) SetExtraData(extraData interface{}) error {
	data, err := json.Marshal(extraData)
	if err != nil {
		return err
	}
	s.ExtraData = data
	return nil
}

// GetExtraData は、追加データを取得します
func (s *Schedule) GetExtraData(v interface{}) error {
	if len(s.ExtraData) == 0 {
		return nil
	}
	return json.Unmarshal(s.ExtraData, v)
}

// IsDeleted は、スケジュールが削除済みかどうかを返します
func (s *Schedule) IsDeleted() bool {
	return !s.DeletedAt.Time.IsZero()
}

// SchedulePerformer は、スケジュールと演者の関連を表します
type SchedulePerformer struct {
	ScheduleID  string    `gorm:"primaryKey" json:"schedule_id"`    // スケジュールID
	PerformerID string    `gorm:"primaryKey" json:"performer_id"`   // 演者ID
	Role        string    `json:"role"`                             // 役割（メインパーソナリティ、アシスタントなど）
	IsGuest     bool      `gorm:"default:false" json:"is_guest"`    // ゲスト出演かどうか
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"` // 作成日時

	// リレーション
	Schedule  *Schedule  `gorm:"foreignKey:ScheduleID" json:"-"`
	Performer *Performer `gorm:"foreignKey:PerformerID" json:"-"`
}

// NewSchedulePerformer は、新しいSchedulePerformerインスタンスを作成します
func NewSchedulePerformer(scheduleID, performerID, role string, isGuest bool) *SchedulePerformer {
	return &SchedulePerformer{
		ScheduleID:  scheduleID,
		PerformerID: performerID,
		Role:        role,
		IsGuest:     isGuest,
	}
}
