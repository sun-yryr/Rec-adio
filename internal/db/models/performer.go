package models

import (
	"time"

	"gorm.io/gorm"
)

// Performer は、番組に出演する人物（声優、パーソナリティなど）の基本情報を表します
type Performer struct {
	ID          string         `gorm:"primaryKey" json:"id"`              // 一意な識別子
	Name        string         `gorm:"not null" json:"name"`              // 名前
	NameReading string         `json:"name_reading"`                      // 名前の読み方（オプション）
	Description string         `json:"description"`                       // 説明（オプション）
	ImageURL    string         `json:"image_url"`                         // 画像URL（オプション）
	OfficialURL string         `json:"official_url"`                      // 公式サイトURL（オプション）
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`  // 作成日時
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`  // 更新日時
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // 削除日時（ソフトデリート用）

	// リレーション
	PlatformAccounts []PlatformAccount `gorm:"foreignKey:PerformerID" json:"-"`
}

// BeforeCreate は、レコード作成前に呼び出されるフック関数です
func (p *Performer) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = GeneratePerformerID()
	}
	return nil
}

// NewPerformer は、新しいPerformerインスタンスを作成します
func NewPerformer(name string) *Performer {
	return &Performer{
		Name: name,
	}
}

// SetNameReading は、名前の読み方を設定します
func (p *Performer) SetNameReading(nameReading string) *Performer {
	p.NameReading = nameReading
	return p
}

// SetDescription は、説明を設定します
func (p *Performer) SetDescription(description string) *Performer {
	p.Description = description
	return p
}

// SetImageURL は、画像URLを設定します
func (p *Performer) SetImageURL(imageURL string) *Performer {
	p.ImageURL = imageURL
	return p
}

// SetOfficialURL は、公式サイトURLを設定します
func (p *Performer) SetOfficialURL(officialURL string) *Performer {
	p.OfficialURL = officialURL
	return p
}
