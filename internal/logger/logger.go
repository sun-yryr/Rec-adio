package logger

import (
	"go.uber.org/zap"

	"github.com/sun-yryr/recoto/internal/config"
)

// NewLogger provides a zap.Logger based on the application environment.
func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	if cfg.Env == "production" {
		return zap.NewProduction()
	}

	return zap.NewDevelopment()
}
