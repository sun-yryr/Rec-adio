package models

import (
	"time"

	"gorm.io/gorm"
)

// ProgramInfo は、番組の基本情報を表します
type ProgramInfo struct {
	ID          string         `gorm:"primaryKey" json:"id"`              // 一意な識別子
	Title       string         `gorm:"not null;index" json:"title"`       // タイトル
	Platform    string         `gorm:"not null;index" json:"platform"`    // プラットフォーム
	Description string         `json:"description"`                       // 説明
	URL         string         `gorm:"not null" json:"url"`               // URL
	ImageURL    string         `json:"image_url"`                         // 画像URL（オプション）
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`  // 作成日時
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`  // 更新日時
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // 削除日時（ソフトデリート用）

	// リレーション
	Performers        []*Performer        `gorm:"many2many:program_performers;" json:"-"`
	PlatformAccounts  []*PlatformAccount  `gorm:"many2many:program_platform_accounts;" json:"-"`
	ProgramPerformers []*ProgramPerformer `gorm:"foreignKey:ProgramInfoID" json:"-"`
}

// BeforeCreate は、レコード作成前に呼び出されるフック関数です
func (p *ProgramInfo) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = GenerateProgramInfoID()
	}
	return nil
}

// NewProgramInfo は、新しいProgramInfoインスタンスを作成します
func NewProgramInfo(title, platform, description, url string) *ProgramInfo {
	return &ProgramInfo{
		Title:       title,
		Platform:    platform,
		Description: description,
		URL:         url,
	}
}

// SetImageURL は、画像URLを設定します
func (p *ProgramInfo) SetImageURL(imageURL string) *ProgramInfo {
	p.ImageURL = imageURL
	return p
}

// ProgramPerformer は、番組と演者の関連を表します
type ProgramPerformer struct {
	ProgramInfoID string    `gorm:"primaryKey" json:"program_info_id"` // 番組情報ID
	PerformerID   string    `gorm:"primaryKey" json:"performer_id"`    // 演者ID
	Role          string    `json:"role"`                              // 役割（メインパーソナリティ、アシスタントなど）
	IsRegular     bool      `gorm:"default:false" json:"is_regular"`   // レギュラー出演者かどうか
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`  // 作成日時

	// リレーション
	ProgramInfo *ProgramInfo `gorm:"foreignKey:ProgramInfoID" json:"-"`
	Performer   *Performer   `gorm:"foreignKey:PerformerID" json:"-"`
}

// NewProgramPerformer は、新しいProgramPerformerインスタンスを作成します
func NewProgramPerformer(programInfoID, performerID, role string, isRegular bool) *ProgramPerformer {
	return &ProgramPerformer{
		ProgramInfoID: programInfoID,
		PerformerID:   performerID,
		Role:          role,
		IsRegular:     isRegular,
	}
}

// ProgramPlatformAccount は、番組とプラットフォームアカウントの関連を表します
type ProgramPlatformAccount struct {
	ProgramInfoID     string    `gorm:"primaryKey" json:"program_info_id"`     // 番組情報ID
	PlatformAccountID string    `gorm:"primaryKey" json:"platform_account_id"` // プラットフォームアカウントID
	Role              string    `json:"role"`                                  // 役割
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`      // 作成日時

	// リレーション
	ProgramInfo     *ProgramInfo     `gorm:"foreignKey:ProgramInfoID" json:"-"`
	PlatformAccount *PlatformAccount `gorm:"foreignKey:PlatformAccountID" json:"-"`
}

// NewProgramPlatformAccount は、新しいProgramPlatformAccountインスタンスを作成します
func NewProgramPlatformAccount(programInfoID, platformAccountID, role string) *ProgramPlatformAccount {
	return &ProgramPlatformAccount{
		ProgramInfoID:     programInfoID,
		PlatformAccountID: platformAccountID,
		Role:              role,
	}
}
