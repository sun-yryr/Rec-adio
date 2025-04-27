package config

import (
	"github.com/caarlos0/env/v11"
)

type logConfig struct {
	Level string `env:"LEVEL" envDefault:"info"`
}

type Config struct {
	Env string    `env:"APP_ENV" envDefault:"development"`
	Log logConfig `envPrefix:"LOG_"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
