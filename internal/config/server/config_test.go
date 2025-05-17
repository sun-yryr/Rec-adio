package server

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cockroachdb/errors"
)

func TestNewDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := NewDefaultConfig()

	if cfg.Log.Level != "info" {
		t.Errorf("Expected Log.Level to be 'info', got %s", cfg.Log.Level)
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("Expected Server.Port to be 8080, got %d", cfg.Server.Port)
	}

	if cfg.Recording.SaveDir != "./data/output" {
		t.Errorf("Expected Recording.SaveDir to be './data/output', got %s", cfg.Recording.SaveDir)
	}
}

func TestLoadConfig_FileNotExists(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.toml")

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Log.Level != "info" {
		t.Errorf("Expected Log.Level to be 'info', got %s", cfg.Log.Level)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}
}

func TestLoadConfig_FileExists(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.toml")
	configContent := `
[log]
level = "error"

[server]
port = 9090

[recording]
save_dir = "/custom/path"
`

	if err := os.WriteFile(configPath, []byte(configContent), 0o600); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Log.Level != "error" {
		t.Errorf("Expected Log.Level to be 'error', got %s", cfg.Log.Level)
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("Expected Server.Port to be 9090, got %d", cfg.Server.Port)
	}

	if cfg.Recording.SaveDir != "/custom/path" {
		t.Errorf("Expected Recording.SaveDir to be '/custom/path', got %s", cfg.Recording.SaveDir)
	}
}

func TestValidateConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setCfg        func(*Config)
		isValid       bool
		expectedError error
	}{
		{
			name:          "valid config",
			setCfg:        func(_ *Config) {},
			isValid:       true,
			expectedError: nil,
		},
		{
			name: "invalid Log.Level",
			setCfg: func(cfg *Config) {
				cfg.Log.Level = "invalid"
			},
			isValid: false,
			expectedError: errors.New(
				"config validation failed: Key: 'Config.Log.Level' Error:Field validation for 'Level' failed on the 'oneof' tag",
			),
		},
		{
			name: "invalid Server.Port",
			setCfg: func(cfg *Config) {
				cfg.Server.Port = 0
			},
			isValid: false,
			expectedError: errors.New(
				"config validation failed: Key: 'Config.Server.Port' Error:Field validation for 'Port' failed on the 'required' tag",
			),
		},
		{
			name: "empty Recording.SaveDir",
			setCfg: func(cfg *Config) {
				cfg.Recording.SaveDir = ""
			},
			isValid: false,
			expectedError: errors.New(
				"config validation failed: Key: 'Config.Recording.SaveDir' Error:Field validation for 'SaveDir' failed on the 'required' tag",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := NewDefaultConfig()
			tt.setCfg(cfg)

			err := ValidateConfig(cfg)
			if tt.isValid {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			} else {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error '%v', got '%v'", tt.expectedError, err)
				}
			}
		})
	}
}
package server

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cockroachdb/errors"
)

func TestNewDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := NewDefaultConfig()

	if cfg.Log.Level != "info" {
		t.Errorf("Expected Log.Level to be 'info', got %s", cfg.Log.Level)
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("Expected Server.Port to be 8080, got %d", cfg.Server.Port)
	}

	if cfg.Recording.SaveDir != "./data/output" {
		t.Errorf("Expected Recording.SaveDir to be './data/output', got %s", cfg.Recording.SaveDir)
	}
}

func TestLoadConfig_FileNotExists(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	// 存在しないファイルパスを指定
	configPath := filepath.Join(tempDir, "config.toml")

	// 設定を読み込む
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// デフォルト値が設定されていることを確認
	if cfg.Log.Level != "info" {
		t.Errorf("Expected Log.Level to be 'info', got %s", cfg.Log.Level)
	}

	// 設定ファイルが作成されていることを確認
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}
}

func TestLoadConfig_FileExists(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	// テスト用の設定ファイルを作成
	configPath := filepath.Join(tempDir, "config.toml")
	configContent := `
[log]
level = "error"

[server]
port = 9090

[recording]
save_dir = "/custom/path"
`

	if err := os.WriteFile(configPath, []byte(configContent), 0o600); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// 設定を読み込む
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Log.Level != "error" {
		t.Errorf("Expected Log.Level to be 'error', got %s", cfg.Log.Level)
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("Expected Server.Port to be 9090, got %d", cfg.Server.Port)
	}

	if cfg.Recording.SaveDir != "/custom/path" {
		t.Errorf("Expected Recording.SaveDir to be '/custom/path', got %s", cfg.Recording.SaveDir)
	}
}

func TestValidateConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setCfg        func(*Config)
		isValid       bool
		expectedError error
	}{
		{
			name:          "valid config",
			setCfg:        func(_ *Config) {},
			isValid:       true,
			expectedError: nil,
		},
		{
			name: "invalid Log.Level",
			setCfg: func(cfg *Config) {
				cfg.Log.Level = "invalid"
			},
			isValid: false,
			expectedError: errors.New(
				"config validation failed: Key: 'Config.Log.Level' Error:Field validation for 'Level' failed on the 'oneof' tag",
			),
		},
		{
			name: "invalid Server.Port",
			setCfg: func(cfg *Config) {
				cfg.Server.Port = 0
			},
			isValid: false,
			expectedError: errors.New(
				"config validation failed: Key: 'Config.Server.Port' Error:Field validation for 'Port' failed on the 'required' tag",
			),
		},
		{
			name: "empty Recording.SaveDir",
			setCfg: func(cfg *Config) {
				cfg.Recording.SaveDir = ""
			},
			isValid: false,
			expectedError: errors.New(
				"config validation failed: Key: 'Config.Recording.SaveDir' Error:Field validation for 'SaveDir' failed on the 'required' tag",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := NewDefaultConfig()
			tt.setCfg(cfg)

			err := ValidateConfig(cfg)
			if tt.isValid {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			} else {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error '%v', got '%v'", tt.expectedError, err)
				}
			}
		})
	}
}

func TestSaveConfig(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	// 設定ファイルのパス
	configPath := filepath.Join(tempDir, "config.toml")

	// テスト用の設定
	cfg := &Config{
		Log: logConfig{
			Level: "debug",
		},
		Server: serverConfig{
			Port: 1234,
		},
		Recording: recordingConfig{
			SaveDir: "/test/path",
		},
	}

	// 設定を保存
	if err := SaveConfig(configPath, cfg); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// 設定ファイルが作成されていることを確認
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// 設定を読み込む
	loadedCfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// 保存した値が正しく読み込まれていることを確認
	if loadedCfg.Log.Level != cfg.Log.Level {
		t.Errorf("Expected Log.Level to be '%s', got '%s'", cfg.Log.Level, loadedCfg.Log.Level)
	}

	if loadedCfg.Server.Port != cfg.Server.Port {
		t.Errorf("Expected Server.Port to be %d, got %d", cfg.Server.Port, loadedCfg.Server.Port)
	}

	if loadedCfg.Recording.SaveDir != cfg.Recording.SaveDir {
		t.Errorf(
			"Expected Recording.SaveDir to be '%s', got '%s'",
			cfg.Recording.SaveDir,
			loadedCfg.Recording.SaveDir,
		)
	}
}
