package recording

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestFinishedEvent(t *testing.T) {
	t.Parallel()

	// 成功ケース
	successEvent := FinishedEvent{
		RecordingID: "test-id",
		Success:     true,
		Error:       nil,
		Timestamp:   time.Now(),
	}

	assert.Equal(t, "test-id", successEvent.RecordingID)
	assert.True(t, successEvent.Success)
	assert.Nil(t, successEvent.Error)
	assert.WithinDuration(t, time.Now(), successEvent.Timestamp, 2*time.Second)

	// 失敗ケース
	errorMessage := "failed to record"
	failureEvent := FinishedEvent{
		RecordingID: "test-id",
		Success:     false,
		Error:       &errorMessage,
		Timestamp:   time.Now(),
	}

	assert.Equal(t, "test-id", failureEvent.RecordingID)
	assert.False(t, failureEvent.Success)
	assert.NotNil(t, failureEvent.Error)
	assert.Equal(t, errorMessage, *failureEvent.Error)
	assert.WithinDuration(t, time.Now(), failureEvent.Timestamp, 2*time.Second)
}

func TestNewFinishedService(t *testing.T) {
	t.Parallel()

	mockBroker := &mockBroker{}
	logger := zaptest.NewLogger(t)

	service := NewFinishedService(mockBroker, logger)

	require.NotNil(t, service)
}
