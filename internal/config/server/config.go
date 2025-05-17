package server

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/cockroachdb/errors"
    "github.com/pelletier/go-toml/v2"
    "github.com/go-playground/validator/v10"
)

type Config struct {
    Log struct {
        Level string `toml:"level" validate:"oneof=debug info warn error"`
    } `toml:"log"`

    Server struct {
        Port int `toml:"port" validate:"required,gt=0"`
    } `toml:"server"`

    Recording struct {
        SaveDir string `toml:"save_dir" validate:"required"`
    } `toml:"recording"`
}

func NewDefaultConfig() *Config {
    cfg := &Config{}
    cfg.Log.Level = "info"
    cfg.Server.Port = 8080
    cfg.Recording.SaveDir = "./data/output"
    return cfg
}

func LoadConfig(path string) (*Config, error) {
    cfg := NewDefaultConfig()

    // Write default config if file does not exist
    if _, err := os.Stat(path); os.IsNotExist(err) {
        dir := filepath.Dir(path)
        if err := os.MkdirAll(dir, 0o755); err != nil {
            return nil, errors.Wrap(err, "failed to create config directory")
        }
        data, err := toml.Marshal(cfg)
        if err != nil {
            return nil, errors.Wrap(err, "failed to marshal default config to TOML")
        }
        if err := os.WriteFile(path, data, 0o600); err != nil {
            return nil, errors.Wrap(err, "failed to write default config file")
        }
        return cfg, nil
    }

    // Load existing config
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, errors.Wrap(err, "failed to read config file")
    }
    if err := toml.Unmarshal(data, cfg); err != nil {
        return nil, errors.Wrap(err, "failed to unmarshal config TOML")
    }

    // Validate loaded config
    if err := ValidateConfig(cfg); err != nil {
        return nil, err
    }
    return cfg, nil
}

func ValidateConfig(cfg *Config) error {
    validate := validator.New()
    if err := validate.Struct(cfg); err != nil {
        return errors.New(fmt.Sprintf("config validation failed: %s", err.Error()))
    }
    return nil
}
package server

import (
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/pelletier/go-toml/v2"
)

type logConfig struct {
	Level string ` + "`toml:\"level\"`" + `
}

type serverConfig struct {
	Port int ` + "`toml:\"port\"`" + `
}

type recordingConfig struct {
	SaveDir string ` + "`toml:\"save_dir\"`" + `
}

type Config struct {
	Log       logConfig       ` + "`toml:\"log\"`" + `
	Server    serverConfig    ` + "`toml:\"server\"`" + `
	Recording recordingConfig ` + "`toml:\"recording\"`" + `
}

func NewDefaultConfig() *Config {
	return &Config{
		Log:       logConfig{Level: "info"},
		Server:    serverConfig{Port: 8080},
		Recording: recordingConfig{SaveDir: "./data/output"},
	}
}

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(home, ".config", "recoto", "daemon.toml")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		cfg := NewDefaultConfig()
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, err
		}
		data, err := toml.Marshal(cfg)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(path, data, 0o600); err != nil {
			return nil, err
		}
		return cfg, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveConfig(path string, cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return err
	}
	return nil
}

func ValidateConfig(cfg *Config) error {
	switch cfg.Log.Level {
	case "debug", "info", "warn", "error":
	default:
		return errors.New("config validation failed: Key: 'Config.Log.Level' Error:Field validation for 'Level' failed on the 'oneof' tag")
	}
	if cfg.Server.Port == 0 {
		return errors.New("config validation failed: Key: 'Config.Server.Port' Error:Field validation for 'Port' failed on the 'required' tag")
	}
	if cfg.Recording.SaveDir == "" {
		return errors.New("config validation failed: Key: 'Config.Recording.SaveDir' Error:Field validation for 'SaveDir' failed on the 'required' tag")
	}
	return nil
}
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
