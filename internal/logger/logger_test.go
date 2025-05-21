package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	config "github.com/sun-yryr/recoto/internal/config/server"
)

func TestNewLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		logLevel  string
		wantError bool
	}{
		{
			name:      "debug level",
			logLevel:  "debug",
			wantError: false,
		},
		{
			name:      "info level",
			logLevel:  "info",
			wantError: false,
		},
		{
			name:      "warn level",
			logLevel:  "warn",
			wantError: false,
		},
		{
			name:      "error level",
			logLevel:  "error",
			wantError: false,
		},
		{
			name:      "invalid level",
			logLevel:  "invalid",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := config.NewDefaultConfig()
			cfg.Log = struct {
				Level string `toml:"level" validate:"required,oneof=debug info warn error"`
			}{
				Level: tt.logLevel,
			}

			logger, err := NewLogger(cfg)

			if tt.wantError {
				require.Error(t, err)
				assert.Nil(t, logger)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, logger)
			}
		})
	}
}

func TestZapConfig(t *testing.T) {
	t.Parallel()

	cfg := zapConfig()

	// 基本的な設定を検証
	assert.Nil(t, cfg.Sampling)
	assert.Contains(t, cfg.OutputPaths, "stdout")
	assert.Contains(t, cfg.ErrorOutputPaths, "stderr")
	assert.Nil(t, cfg.Sampling)
	// DurationEncoder関数を直接比較するのは難しいので、型の比較だけを行う
	assert.NotNil(t, cfg.EncoderConfig.EncodeDuration)
}
