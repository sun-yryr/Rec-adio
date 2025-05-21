package recorder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sun-yryr/recoto/internal/config/radiko"
	"go.uber.org/zap/zaptest"
)

func TestNewRadikoRecorderFromConfig(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	config := &radiko.Config{
		URL:    "http://radiko.jp/v3/program/today/JP27.xml",
		AreaID: "JP27",
	}

	recorder := NewRadikoRecorderFromConfig(logger, config)

	assert.NotNil(t, recorder)
	assert.Equal(t, "RadikoRecorder", recorder.GetName())
	assert.Equal(t, config.URL, recorder.radikoURL)
	assert.Equal(t, config.AreaID, recorder.areaID)
}
