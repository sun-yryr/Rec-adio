// Package config は、アプリケーションの設定を管理します.
package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/cockroachdb/errors"
)

type logConfig struct {
	Level string `env:"LEVEL" envDefault:"info"`
}

type serverConfig struct {
	Port int `env:"PORT" envDefault:"8080"`
}

// Config represents the application configuration.
type Config struct {
	Env    string       `env:"APP_ENV" envDefault:"development"`
	Log    logConfig    `                                       envPrefix:"LOG_"`
	Server serverConfig `                                       envPrefix:"SERVER_"`
}

// Load loads the application configuration from environment variables.
func Load() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, errors.Wrap(err, "failed to parse config")
	}

	return &cfg, nil
}
