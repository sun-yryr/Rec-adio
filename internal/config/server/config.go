// Package server は、アプリケーションの設定を管理する
package server

import (
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/go-playground/validator/v10"
	"github.com/pelletier/go-toml/v2"

	"github.com/sun-yryr/recoto/internal/fileutil"
)

type logConfig struct {
	Level string `toml:"level" validate:"required,oneof=debug info warn error"`
}

type serverConfig struct {
	Port int `toml:"port" validate:"required,gt=0,lte=65535"`
}

type recordingConfig struct {
	SaveDir string `toml:"save_dir" validate:"required"`
}

// Config はアプリケーションの設定を表す。
type Config struct {
	Log       logConfig       `toml:"log"       validate:"required"`
	Server    serverConfig    `toml:"server"    validate:"required"`
	Recording recordingConfig `toml:"recording" validate:"required"`
}

// NewDefaultConfig はデフォルト値を持つ設定を作成する。
func NewDefaultConfig() *Config {
	return &Config{
		Log: logConfig{
			Level: "info",
		},
		Server: serverConfig{
			Port: 8080, //nolint:mnd
		},
		Recording: recordingConfig{
			SaveDir: "./data/output",
		},
	}
}

// GetDefaultConfigPath はデフォルトの設定ファイルパスを返す。
func GetDefaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get home directory")
	}

	return filepath.Join(home, ".config", "recoto", "daemon.toml"), nil
}

// LoadConfig は設定ファイルから設定を読み込む。存在しない場合はデフォルト設定を作成する。
func LoadConfig(configPath string) (*Config, error) {
	var err error

	// 設定ファイルのパスを決定
	if configPath == "" {
		configPath, err = GetDefaultConfigPath()
		if err != nil {
			return nil, err
		}
	}

	config := NewDefaultConfig()

	// ファイルが存在するか確認
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		// ディレクトリが存在することを確認
		dir := filepath.Dir(configPath)
		if err := os.MkdirAll(dir, fileutil.DefaultDirPerm); err != nil {
			return nil, errors.Wrap(err, "failed to create config directory")
		}

		if err := SaveConfig(configPath, config); err != nil {
			return nil, err
		}

		return config, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to check config file")
	}

	data, err := os.ReadFile(configPath) //nolint:gosec
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}

	if err := toml.Unmarshal(data, config); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}

	return config, nil
}

// SaveConfig は設定をファイルに保存する。
func SaveConfig(configPath string, config *Config) error {
	if err := ValidateConfig(config); err != nil {
		return errors.Wrap(err, "failed to validate config")
	}

	data, err := toml.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "failed to marshal config")
	}

	if err := os.WriteFile(configPath, data, fileutil.DefaultFilePerm); err != nil {
		return errors.Wrap(err, "failed to write config file")
	}

	return nil
}

// ValidateConfig は設定値のバリデーションを行う。
func ValidateConfig(config *Config) error {
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return errors.Wrap(err, "config validation failed")
	}

	return nil
}
