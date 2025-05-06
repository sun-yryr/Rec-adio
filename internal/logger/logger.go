// Package logger は、アプリケーションロガーを提供する
package logger

import (
	"github.com/cockroachdb/errors"
	"go.uber.org/zap"

	"github.com/sun-yryr/recoto/internal/config"
)

// NewLogger は、アプリケーション環境に基づいてzap.Loggerを作成する。
func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	if cfg.Env == "production" {
		logger, err := zap.NewProduction()
		if err != nil {
			return nil, errors.Wrap(err, "failed to create production logger")
		}

		return logger, nil
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create development logger")
	}

	return logger, nil
}
