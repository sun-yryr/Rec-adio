package server

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := NewDefaultConfig()

	assert.Equal(t, "info", cfg.Log.Level)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "./data/output", cfg.Recording.SaveDir)
}

func TestGetDefaultConfigPath(t *testing.T) {
	t.Parallel()

	path, err := GetDefaultConfigPath()
	require.NoError(t, err)

	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	expectedPath := filepath.Join(homeDir, ".config", "recoto", "daemon.toml")
	assert.Equal(t, expectedPath, path)
}

func TestLoadConfig_FileNotExists(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	// 存在しないファイルパスを指定
	configPath := filepath.Join(tempDir, "config.toml")

	// 設定を読み込む
	cfg, err := LoadConfig(configPath)
	require.NoError(t, err)

	// デフォルト値が設定されていることを確認
	assert.Equal(t, "info", cfg.Log.Level)
	assert.FileExists(t, configPath)
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

	require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0o600))

	// 設定を読み込む
	cfg, err := LoadConfig(configPath)
	require.NoError(t, err)

	assert.Equal(t, "error", cfg.Log.Level)
	assert.Equal(t, 9090, cfg.Server.Port)
	assert.Equal(t, "/custom/path", cfg.Recording.SaveDir)
}

func TestLoadConfig_DefaultPath(t *testing.T) {
	t.Parallel()

	// 一時的に環境変数をモックして HOME を設定
	oldHome := os.Getenv("HOME")

	t.Cleanup(func() {
		t.Setenv("HOME", oldHome)
	})

	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	// デフォルトパスに設定ファイルを作成
	configDir := filepath.Join(tempDir, ".config", "recoto")
	require.NoError(t, os.MkdirAll(configDir, 0o755))

	configPath := filepath.Join(configDir, "daemon.toml")
	configContent := `
[log]
level = "debug"

[server]
port = 7070

[recording]
save_dir = "/test/dir"
`
	require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0o600))

	// 空のパスを指定してデフォルトパスから読み込む
	cfg, err := LoadConfig("")
	require.NoError(t, err)

	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, 7070, cfg.Server.Port)
	assert.Equal(t, "/test/dir", cfg.Recording.SaveDir)
}

func TestLoadConfig_InvalidTOML(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	// 不正なTOML形式の設定ファイルを作成
	configPath := filepath.Join(tempDir, "config.toml")
	invalidContent := `
[log]
level = "error"

[server]
port = "not a number" # 整数であるべき
`

	require.NoError(t, os.WriteFile(configPath, []byte(invalidContent), 0o600))

	// 設定を読み込む
	_, err := LoadConfig(configPath)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal config")
}

func TestLoadConfig_StatError(t *testing.T) {
	t.Parallel()

	// ファイルのチェックが失敗するケース
	// これはOSによって振る舞いが異なるため、全環境で動作するようにはしにくい
	// 例えば、アクセス権のないディレクトリ内のファイルなど
	configPath := "/root/.impossible/config.toml" // 通常のユーザーではアクセスできないパス

	// 設定を読み込む試行
	_, err := LoadConfig(configPath)
	// エラーが発生することだけ確認
	require.Error(t, err)
}

func TestValidateConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setCfg        func(*Config)
		isValid       bool
		expectedField string
		expectedTag   string
	}{
		{
			name:          "valid config",
			setCfg:        func(_ *Config) {},
			isValid:       true,
			expectedField: "",
			expectedTag:   "",
		},
		{
			name: "invalid Log.Level",
			setCfg: func(cfg *Config) {
				cfg.Log.Level = "invalid"
			},
			isValid:       false,
			expectedField: "Level",
			expectedTag:   "oneof",
		},
		{
			name: "invalid Server.Port",
			setCfg: func(cfg *Config) {
				cfg.Server.Port = 0
			},
			isValid:       false,
			expectedField: "Port",
			expectedTag:   "required",
		},
		{
			name: "empty Recording.SaveDir",
			setCfg: func(cfg *Config) {
				cfg.Recording.SaveDir = ""
			},
			isValid:       false,
			expectedField: "SaveDir",
			expectedTag:   "required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := NewDefaultConfig()
			tt.setCfg(cfg)

			err := ValidateConfig(cfg)

			// 正しい設定を想定する場合はNoErrorを確認して終了
			if tt.isValid {
				assert.NoError(t, err)

				return
			}

			var validationErrors validator.ValidationErrors

			require.ErrorAs(t, err, &validationErrors)

			assert.NotEmpty(t, validationErrors)

			// エラーのフィールド名とタグを検証
			found := false

			for _, fieldErr := range validationErrors {
				if fieldErr.Field() == tt.expectedField && fieldErr.Tag() == tt.expectedTag {
					found = true

					break
				}
			}

			if !found {
				t.Errorf("Expected validation error for field '%s' with tag '%s', got %v",
					tt.expectedField, tt.expectedTag, validationErrors)
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
	require.NoError(t, SaveConfig(configPath, cfg))
	// 設定ファイルが作成されていることを確認
	require.FileExists(t, configPath)

	// 設定を読み込む
	loadedCfg, err := LoadConfig(configPath)
	require.NoError(t, err)

	// 保存した値が正しく読み込まれていることを確認
	assert.Equal(t, cfg.Log.Level, loadedCfg.Log.Level)
	assert.Equal(t, cfg.Server.Port, loadedCfg.Server.Port)
	assert.Equal(t, cfg.Recording.SaveDir, loadedCfg.Recording.SaveDir)
}

func TestSaveConfig_InvalidConfig(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.toml")

	// 無効な設定
	invalidCfg := &Config{
		Log: logConfig{
			Level: "invalid", // "debug", "info", "warn", "error" のいずれかでなければならない
		},
		Server: serverConfig{
			Port: 8080,
		},
		Recording: recordingConfig{
			SaveDir: "/test/path",
		},
	}

	// 無効な設定で保存を試みる
	err := SaveConfig(configPath, invalidCfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to validate config")
}

func TestSaveConfig_WriteError(t *testing.T) {
	t.Parallel()

	// 書き込み権限のないパス
	configPath := "/root/impossible_config.toml"

	cfg := NewDefaultConfig()

	// 書き込み権限のないパスに保存を試みる
	err := SaveConfig(configPath, cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to write config file")
}
