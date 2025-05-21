// Package radiko は、Radikoの設定を管理する
package radiko

import (
	"os"

	"github.com/cockroachdb/errors"
	"github.com/go-playground/validator/v10"
	"github.com/pelletier/go-toml/v2"
)

// Config はRadikoの設定を表す。
type Config struct {
	URL    string `toml:"url"     validate:"required,url"`
	AreaID string `toml:"area_id" validate:"required"`
}

// NewDefaultConfig はデフォルト値を持つRadiko設定を作成する。
func NewDefaultConfig() *Config {
	return &Config{
		URL:    "http://radiko.jp/v3/program/today/JP13.xml",
		AreaID: "JP13",
	}
}

// LoadConfig はTOMLデータからRadiko設定を読み込む。
func LoadConfig(data []byte) (*Config, error) {
	config := NewDefaultConfig()

	// radikoセクションが存在しない場合はデフォルト設定を返す
	var tmp struct {
		Radiko *Config `toml:"radiko"`
	}

	if err := toml.Unmarshal(data, &tmp); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal radiko config")
	}

	if tmp.Radiko != nil {
		config = tmp.Radiko
	}

	if err := ValidateConfig(config); err != nil {
		return nil, errors.Wrap(err, "failed to validate radiko config")
	}

	return config, nil
}

// LoadConfigFromFile はファイルからRadiko設定を読み込む。
func LoadConfigFromFile(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath) //nolint:gosec
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}

	return LoadConfig(data)
}

// ValidateConfig はRadiko設定値のバリデーションを行う。
func ValidateConfig(config *Config) error {
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return errors.Wrap(err, "radiko config validation failed")
	}

	return nil
}
