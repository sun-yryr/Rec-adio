package recorder

import (
	"go.uber.org/zap"

	"github.com/sun-yryr/recoto/internal/config/radiko"
)

// NewRadikoRecorderFromConfig はRadiko設定からRadikoRecorderを生成する。
func NewRadikoRecorderFromConfig(logger *zap.Logger, config *radiko.Config) *RadikoRecorder {
	return NewRadikoRecorder(logger, RadikoConfig{
		RadikoURL: config.URL,
		AreaID:    config.AreaID,
	})
}
