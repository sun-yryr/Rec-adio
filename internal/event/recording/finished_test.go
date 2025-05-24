package recording

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFinishedEvent(t *testing.T) {
	t.Parallel()

	// FinishedEventの基本構造が期待通りか確認
	event := FinishedEvent{
		RecordingID: "test-id",
		Success:     true,
		Error:       nil,
		Timestamp:   time.Now(),
	}

	assert.Equal(t, "test-id", event.RecordingID)
	assert.True(t, event.Success)
	assert.Nil(t, event.Error)
	assert.WithinDuration(t, time.Now(), event.Timestamp, 2*time.Second)

	// エラーケースのテスト
	errorMsg := "test error"
	errorEvent := FinishedEvent{
		RecordingID: "test-id-error",
		Success:     false,
		Error:       &errorMsg,
		Timestamp:   time.Now(),
	}

	assert.Equal(t, "test-id-error", errorEvent.RecordingID)
	assert.False(t, errorEvent.Success)
	assert.NotNil(t, errorEvent.Error)
	assert.Equal(t, "test error", *errorEvent.Error)
}

func TestNewFinishedService(t *testing.T) {
	t.Parallel()

	mockBroker := &mockBroker{}

	service := NewFinishedService(mockBroker)

	require.NotNil(t, service)
}
