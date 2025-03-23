package models

import (
	"time"

	"gorm.io/gorm"
)

// PerformerSubscription は、演者の購読情報を表します
type PerformerSubscription struct {
	ID                  string         `gorm:"primaryKey" json:"id"`                              // 一意な識別子
	PerformerID         string         `gorm:"not null;index" json:"performer_id"`                // 演者ID
	Enabled             bool           `gorm:"not null;default:true" json:"enabled"`              // 有効かどうか
	NotificationEnabled bool           `gorm:"not null;default:true" json:"notification_enabled"` // 通知を有効にするかどうか
	CreatedAt           time.Time      `gorm:"autoCreateTime" json:"created_at"`                  // 作成日時
	UpdatedAt           time.Time      `gorm:"autoUpdateTime" json:"updated_at"`                  // 更新日時
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`                 // 削除日時（ソフトデリート用）
	// TODO: 論理削除いらない。消す

	// リレーション
	Performer *Performer `gorm:"foreignKey:PerformerID" json:"-"`
}

// BeforeCreate は、レコード作成前に呼び出されるフック関数です
func (s *PerformerSubscription) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = GeneratePerformerSubscriptionID()
	}
	return nil
}

// NewPerformerSubscription は、新しいPerformerSubscriptionインスタンスを作成します
func NewPerformerSubscription(performerID string) *PerformerSubscription {
	return &PerformerSubscription{
		PerformerID:         performerID,
		Enabled:             true,
		NotificationEnabled: true,
	}
}

// SetEnabled は、有効かどうかを設定します
func (s *PerformerSubscription) SetEnabled(enabled bool) *PerformerSubscription {
	s.Enabled = enabled
	return s
}

// SetNotificationEnabled は、通知を有効にするかどうかを設定します
func (s *PerformerSubscription) SetNotificationEnabled(enabled bool) *PerformerSubscription {
	s.NotificationEnabled = enabled
	return s
}

// PlatformAccountSubscription は、プラットフォームアカウントの購読情報を表します
type PlatformAccountSubscription struct {
	ID                  string         `gorm:"primaryKey" json:"id"`                              // 一意な識別子
	PlatformAccountID   string         `gorm:"not null;index" json:"platform_account_id"`         // プラットフォームアカウントID
	Enabled             bool           `gorm:"not null;default:true" json:"enabled"`              // 有効かどうか
	NotificationEnabled bool           `gorm:"not null;default:true" json:"notification_enabled"` // 通知を有効にするかどうか
	CreatedAt           time.Time      `gorm:"autoCreateTime" json:"created_at"`                  // 作成日時
	UpdatedAt           time.Time      `gorm:"autoUpdateTime" json:"updated_at"`                  // 更新日時
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`                 // 削除日時（ソフトデリート用）

	// リレーション
	PlatformAccount *PlatformAccount `gorm:"foreignKey:PlatformAccountID" json:"-"`
}

// BeforeCreate は、レコード作成前に呼び出されるフック関数です
func (s *PlatformAccountSubscription) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = GeneratePlatformAccountSubscriptionID()
	}
	return nil
}

// NewPlatformAccountSubscription は、新しいPlatformAccountSubscriptionインスタンスを作成します
func NewPlatformAccountSubscription(platformAccountID string) *PlatformAccountSubscription {
	return &PlatformAccountSubscription{
		PlatformAccountID:   platformAccountID,
		Enabled:             true,
		NotificationEnabled: true,
	}
}

// SetEnabled は、有効かどうかを設定します
func (s *PlatformAccountSubscription) SetEnabled(enabled bool) *PlatformAccountSubscription {
	s.Enabled = enabled
	return s
}

// SetNotificationEnabled は、通知を有効にするかどうかを設定します
func (s *PlatformAccountSubscription) SetNotificationEnabled(enabled bool) *PlatformAccountSubscription {
	s.NotificationEnabled = enabled
	return s
}
