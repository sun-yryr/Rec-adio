package models

import (
	"time"

	"gorm.io/gorm"
)

// PlatformAccount は、各プラットフォーム上のアカウント情報を表します
type PlatformAccount struct {
	ID          string         `gorm:"primaryKey" json:"id"`              // 一意な識別子
	PerformerID string         `gorm:"index" json:"performer_id"`         // 関連する演者のID（オプション）
	Platform    string         `gorm:"not null;index" json:"platform"`    // プラットフォーム名（twitter, bilibili など）
	Username    string         `gorm:"not null" json:"username"`          // ユーザー名
	DisplayName string         `json:"display_name"`                      // 表示名（オプション）
	AvatarURL   string         `json:"avatar_url"`                        // アバターURL（オプション）
	IsOfficial  bool           `gorm:"default:false" json:"is_official"`  // 公式アカウントかどうか
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`  // 作成日時
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`  // 更新日時
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // 削除日時（ソフトデリート用）

	// リレーション
	Performer *Performer `gorm:"foreignKey:PerformerID" json:"-"`
}

// BeforeCreate は、レコード作成前に呼び出されるフック関数です
func (a *PlatformAccount) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = GeneratePlatformAccountID()
	}
	return nil
}

// NewPlatformAccount は、新しいPlatformAccountインスタンスを作成します
func NewPlatformAccount(platform, username string) *PlatformAccount {
	return &PlatformAccount{
		Platform:   platform,
		Username:   username,
		IsOfficial: false,
	}
}

// SetPerformerID は、関連する演者のIDを設定します
func (a *PlatformAccount) SetPerformerID(performerID string) *PlatformAccount {
	a.PerformerID = performerID
	return a
}

// SetDisplayName は、表示名を設定します
func (a *PlatformAccount) SetDisplayName(displayName string) *PlatformAccount {
	a.DisplayName = displayName
	return a
}

// SetAvatarURL は、アバターURLを設定します
func (a *PlatformAccount) SetAvatarURL(avatarURL string) *PlatformAccount {
	a.AvatarURL = avatarURL
	return a
}

// SetIsOfficial は、公式アカウントかどうかを設定します
func (a *PlatformAccount) SetIsOfficial(isOfficial bool) *PlatformAccount {
	a.IsOfficial = isOfficial
	return a
}
