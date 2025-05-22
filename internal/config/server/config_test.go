package server

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sun-yryr/recoto/internal/fileutil"
)

type statErrorFS struct{ fileutil.FileSystem }

func (statErrorFS) Stat(string) (os.FileInfo, error) { return nil, errStat }

type writeErrorFS struct{ fileutil.FileSystem }

func (writeErrorFS) WriteFile(string, []byte, os.FileMode) error { return errWrite }

var (
	errStat  = errors.New("stat error")
	errWrite = errors.New("write error")
)

const testConfigPath = "/config.toml"

func TestNewDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := NewDefaultConfig()

	assert.Equal(t, "info", cfg.Log.Level)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "./data/output", cfg.Recording.SaveDir)
}

func TestGetDefaultConfigPath(t *testing.T) {
	fs := fileutil.NewMemFS()

	t.Setenv("HOME", "/home/test")

	path, err := GetDefaultConfigPathFS(fs)
	require.NoError(t, err)

	expectedPath := filepath.Join("/home/test", ".config", "recoto", "daemon.toml")
	assert.Equal(t, expectedPath, path)
}

func TestLoadConfig_FileNotExists(t *testing.T) {
	t.Parallel()

	fs := fileutil.NewMemFS()
	cfg, err := LoadConfigFS(fs, testConfigPath)
	require.NoError(t, err)

	assert.Equal(t, "info", cfg.Log.Level)

	_, statErr := fs.Stat(testConfigPath)
	assert.NoError(t, statErr)
}

func TestLoadConfig_FileExists(t *testing.T) {
	t.Parallel()

	fs := fileutil.NewMemFS()
	configPath := testConfigPath
	configContent := `
[log]
level = "error"

[server]
port = 9090

[recording]
save_dir = "/custom/path"
`

	require.NoError(t, fs.WriteFile(configPath, []byte(configContent), 0o600))
	cfg, err := LoadConfigFS(fs, configPath)
	require.NoError(t, err)

	assert.Equal(t, "error", cfg.Log.Level)
	assert.Equal(t, 9090, cfg.Server.Port)
	assert.Equal(t, "/custom/path", cfg.Recording.SaveDir)
}

func TestLoadConfig_DefaultPath(t *testing.T) {
	// 一時的に環境変数をモックして HOME を設定
	oldHome := os.Getenv("HOME")

	t.Cleanup(func() {
		t.Setenv("HOME", oldHome)
	})

	tempDir := t.TempDir()
	fs := fileutil.NewMemFS()

	t.Setenv("HOME", tempDir)
	// デフォルトパスに設定ファイルを作成
	configDir := filepath.Join(tempDir, ".config", "recoto")
	require.NoError(t, fs.MkdirAll(configDir, 0o755))

	configPath := filepath.Join(configDir, "daemon.toml")
	configContent := `
[log]
level = "debug"

[server]
port = 7070

[recording]
save_dir = "/test/dir"
`
	require.NoError(t, fs.WriteFile(configPath, []byte(configContent), 0o600))

	cfg, err := LoadConfigFS(fs, "")
	require.NoError(t, err)

	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, 7070, cfg.Server.Port)
	assert.Equal(t, "/test/dir", cfg.Recording.SaveDir)
}

func TestLoadConfig_InvalidTOML(t *testing.T) {
	t.Parallel()

	fs := fileutil.NewMemFS()
	configPath := testConfigPath
	invalidContent := `
[log]
level = "error"

[server]
port = "not a number" # 整数であるべき
`

	require.NoError(t, fs.WriteFile(configPath, []byte(invalidContent), 0o600))

	_, err := LoadConfigFS(fs, configPath)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal config")
}

func TestLoadConfig_StatError(t *testing.T) {
	t.Parallel()

	fs := statErrorFS{fileutil.NewMemFS()}
	_, err := LoadConfigFS(fs, testConfigPath)
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

	fs := fileutil.NewMemFS()
	configPath := testConfigPath

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

	require.NoError(t, SaveConfigFS(fs, configPath, cfg))
	_, statErr := fs.Stat(configPath)
	require.NoError(t, statErr)

	loadedCfg, err := LoadConfigFS(fs, configPath)
	require.NoError(t, err)

	// 保存した値が正しく読み込まれていることを確認
	assert.Equal(t, cfg.Log.Level, loadedCfg.Log.Level)
	assert.Equal(t, cfg.Server.Port, loadedCfg.Server.Port)
	assert.Equal(t, cfg.Recording.SaveDir, loadedCfg.Recording.SaveDir)
}

func TestSaveConfig_InvalidConfig(t *testing.T) {
	t.Parallel()

	fs := fileutil.NewMemFS()
	configPath := testConfigPath

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
	err := SaveConfigFS(fs, configPath, invalidCfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to validate config")
}

func TestSaveConfig_WriteError(t *testing.T) {
	t.Parallel()

	mem := fileutil.NewMemFS()
	ro := writeErrorFS{mem}
	cfg := NewDefaultConfig()

	err := SaveConfigFS(ro, testConfigPath, cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to write config file")
}
