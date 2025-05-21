package radiko

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultConfig(t *testing.T) {
	t.Parallel()

	config := NewDefaultConfig()
	assert.Equal(t, "http://radiko.jp/v3/program/today/JP13.xml", config.URL)
	assert.Equal(t, "JP13", config.AreaID)
}

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		toml     string
		expected *Config
		wantErr  bool
	}{
		{
			name: "valid config",
			toml: `
[radiko]
url = "http://radiko.jp/v3/program/today/JP13.xml"
area_id = "JP13"
`,
			expected: &Config{
				URL:    "http://radiko.jp/v3/program/today/JP13.xml",
				AreaID: "JP13",
			},
			wantErr: false,
		},
		{
			name: "custom config",
			toml: `
[radiko]
url = "http://radiko.jp/v3/program/today/JP27.xml"
area_id = "JP27"
`,
			expected: &Config{
				URL:    "http://radiko.jp/v3/program/today/JP27.xml",
				AreaID: "JP27",
			},
			wantErr: false,
		},
		{
			name: "missing radiko section",
			toml: `
[other]
key = "value"
`,
			expected: NewDefaultConfig(),
			wantErr:  false,
		},
		{
			name: "invalid url",
			toml: `
[radiko]
url = "invalid-url"
area_id = "JP13"
`,
			expected: nil,
			wantErr:  true,
		},
		{
			name: "missing area_id",
			toml: `
[radiko]
url = "http://radiko.jp/v3/program/today/JP13.xml"
`,
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			config, err := LoadConfig([]byte(tt.toml))
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, config)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, config)
			}
		})
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	t.Parallel()

	// テスト用の一時ファイルを作成
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")

	toml := `
[radiko]
url = "http://radiko.jp/v3/program/today/JP27.xml"
area_id = "JP27"
`

	err := os.WriteFile(configPath, []byte(toml), 0o600)
	require.NoError(t, err)

	config, err := LoadConfigFromFile(configPath)
	require.NoError(t, err)
	assert.Equal(t, "http://radiko.jp/v3/program/today/JP27.xml", config.URL)
	assert.Equal(t, "JP27", config.AreaID)

	// 存在しないファイル
	_, err = LoadConfigFromFile(filepath.Join(dir, "non-existent.toml"))
	assert.Error(t, err)
}

func TestValidateConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				URL:    "http://radiko.jp/v3/program/today/JP13.xml",
				AreaID: "JP13",
			},
			wantErr: false,
		},
		{
			name: "invalid url",
			config: &Config{
				URL:    "invalid-url",
				AreaID: "JP13",
			},
			wantErr: true,
		},
		{
			name: "empty area_id",
			config: &Config{
				URL:    "http://radiko.jp/v3/program/today/JP13.xml",
				AreaID: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateConfig(tt.config)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
