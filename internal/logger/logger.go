package logger

import (
	"github.com/sun-yryr/recoto/internal/config"
	"go.uber.org/zap"
)

// NewLogger provides a zap.Logger based on the application environment.
func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	if cfg.Env == "production" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}
