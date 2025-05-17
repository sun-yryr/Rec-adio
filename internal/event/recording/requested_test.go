package recording

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/sun-yryr/recoto/internal/broker"
	"github.com/sun-yryr/recoto/internal/domain"
)

func TestNewRequestedEvent(t *testing.T) {
	t.Parallel()

	// テスト用のSource作成
	source, err := domain.NewURLSource("http://example.com/stream")
	require.NoError(t, err)

	// テスト用のRecording作成
	recording, err := domain.NewRecording(source, "/path/to/output.mp3", 30*time.Minute)
	require.NoError(t, err)

	// テスト関数実行
	event := NewRequestedEvent(*recording)

	// 結果の検証
	assert.Equal(t, string(recording.ID), event.RecordingID)
	assert.Equal(t, recording.Status, event.Status)
	assert.Equal(t, recording.Source, event.Source)
	assert.Equal(t, recording.Output, event.Output)
	assert.Equal(t, recording.Duration, event.Duration)
	// タイムスタンプは現在時刻に近いことを確認
	assert.WithinDuration(t, time.Now(), event.Timestamp, 2*time.Second)
}

func TestRequestedEvent_ToDomain(t *testing.T) {
	t.Parallel()

	// テスト用のSource作成
	source, err := domain.NewURLSource("http://example.com/stream")
	require.NoError(t, err)

	// テスト用のRequestedEvent作成
	event := &RequestedEvent{
		RecordingID: "test-id",
		Status:      domain.RecordingStatusRequested,
		Source:      source,
		Output:      "/path/to/output.mp3",
		Duration:    30 * time.Minute,
		Timestamp:   time.Now(),
	}

	// テスト関数実行
	recording := event.ToDomain()

	// 結果の検証
	assert.Equal(t, domain.RecordingID(event.RecordingID), recording.ID)
	assert.Equal(t, event.Status, recording.Status)
	assert.Equal(t, event.Source, recording.Source)
	assert.Equal(t, event.Output, recording.Output)
	assert.Equal(t, event.Duration, recording.Duration)
}

func TestNewRequestedService(t *testing.T) {
	t.Parallel()

	mockBroker := &mockBroker{}
	logger := zaptest.NewLogger(t)

	service := NewRequestedService(mockBroker, logger)

	require.NotNil(t, service)
}

// mockBroker はテスト用のブローカーモック
type mockBroker struct{}

func (m *mockBroker) Publish(_ context.Context, _ string, _ []byte) error {
	return nil
}

func (m *mockBroker) Subscribe(
	_ context.Context,
	_ string,
	_ func(message []byte),
) (broker.UnsubscribeFunc, error) {
	return func() error { return nil }, nil
}

func (m *mockBroker) Close() error {
	return nil
}