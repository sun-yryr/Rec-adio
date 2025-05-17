// Package logger は、アプリケーションロガーを提供する
package logger

import (
	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	config "github.com/sun-yryr/recoto/internal/config/server"
)

// NewLogger は、アプリケーション環境に基づいてzap.Loggerを作成する。
func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	zapConfig := zapConfig()

	level, err := zapcore.ParseLevel(cfg.Log.Level)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse log level")
	}

	zapConfig.Level.SetLevel(level)

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create logger")
	}

	return logger, nil
}

func zapConfig() *zap.Config {
	cfg := zap.NewProductionConfig()

	cfg.Sampling = nil
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder

	return &cfg
}
