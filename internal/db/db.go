package db

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Config は、データベースの設定を表す構造体です
type Config struct {
	Path string // データベースファイルのパス
}

// New は、新しいDBインスタンスを作成します
func New(config Config) (*gorm.DB, error) {
	// データベースファイルのディレクトリが存在しない場合は作成
	dbDir := filepath.Dir(config.Path)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	// データベース接続を開く
	db, err := gorm.Open(sqlite.Open(config.Path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	DB = db

	return db, nil
}
