package recording

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartedEvent(t *testing.T) {
	t.Parallel()

	// StartedEventの基本構造が期待通りか確認
	event := StartedEvent{
		RecordingID: "test-id",
		Timestamp:   time.Now(),
	}

	assert.Equal(t, "test-id", event.RecordingID)
	assert.WithinDuration(t, time.Now(), event.Timestamp, 2*time.Second)
}

func TestNewStartedService(t *testing.T) {
	t.Parallel()

	mockBroker := &mockBroker{}

	service := NewStartedService(mockBroker)

	require.NotNil(t, service)
}
