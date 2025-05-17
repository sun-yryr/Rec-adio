package server

import (
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/pelletier/go-toml/v2"
	"github.com/go-playground/validator/v10"
)

type LogConfig struct {
	Level string `toml:"level" validate:"oneof=debug info warn error"`
}

type ServerConfig struct {
	Port int `toml:"port" validate:"required"`
}

type RecordingConfig struct {
	SaveDir string `toml:"save_dir" validate:"required"`
}

type Config struct {
	Log       LogConfig       `toml:"log"`
	Server    ServerConfig    `toml:"server"`
	Recording RecordingConfig `toml:"recording"`
}

func NewDefaultConfig() *Config {
	return &Config{
		Log:       LogConfig{Level: "info"},
		Server:    ServerConfig{Port: 8080},
		Recording: RecordingConfig{SaveDir: "./data/output"},
	}
}

func ValidateConfig(cfg *Config) error {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return errors.New("config validation failed: " + err.Error())
	}
	return nil
}

func LoadConfig(path string) (*Config, error) {
	cfg := NewDefaultConfig()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, errors.Wrap(err, "failed to create config directory")
		}

		data, err := toml.Marshal(cfg)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal default config")
		}

		if err := os.WriteFile(path, data, 0o600); err != nil {
			return nil, errors.Wrap(err, "failed to write default config file")
		}
		return cfg, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to stat config file")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}

	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config file")
	}

	if err := ValidateConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}