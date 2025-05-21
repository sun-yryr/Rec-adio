package grpcserver

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/sun-yryr/recoto/internal/broker"
	config "github.com/sun-yryr/recoto/internal/config/server"
	"github.com/sun-yryr/recoto/internal/event/recording"
	"github.com/sun-yryr/recoto/internal/eventutil"
	recordingv1 "github.com/sun-yryr/recoto/pkg/api/recoto/recording/v1"
)

// mockRecordingBroker はテスト用のブローカーモック。
type mockRecordingBroker struct {
	publishFunc func(ctx context.Context, subject string, message []byte) error
}

func (m *mockRecordingBroker) Publish(ctx context.Context, subject string, message []byte) error {
	if m.publishFunc != nil {
		return m.publishFunc(ctx, subject, message)
	}

	return nil
}

func (m *mockRecordingBroker) Subscribe(
	_ context.Context,
	_ string,
	_ func(message []byte),
) (broker.UnsubscribeFunc, error) {
	return func() error { return nil }, nil
}

func (m *mockRecordingBroker) Close() error {
	return nil
}

func TestNewRecordingService(t *testing.T) {
	t.Parallel()

	// 設定とイベントサービスの準備
	cfg := &config.Config{
		Recording: struct {
			SaveDir string `toml:"save_dir" validate:"required"`
		}{
			SaveDir: "/test/save/dir",
		},
	}
	logger := zaptest.NewLogger(t)
	mockBroker := &mockRecordingBroker{}
	service := eventutil.NewEventService[recording.RequestedEvent](
		mockBroker,
		logger,
		"test.subject",
	)

	// サービスの作成
	recordingService := NewRecordingService(cfg, service)

	// 結果の検証
	assert.NotNil(t, recordingService)
	assert.Equal(t, service, recordingService.requestedService)
	assert.Equal(t, "/test/save/dir", recordingService.saveDir)
}

func TestRecordingService_StartFromURL_Success(t *testing.T) {
	t.Parallel()

	// 一時ディレクトリの作成
	tempDir := t.TempDir()

	// publishされたイベントをキャプチャするためのモック
	var (
		capturedSubject string
		capturedMessage []byte
	)

	mockBroker := &mockRecordingBroker{
		publishFunc: func(_ context.Context, subject string, message []byte) error {
			capturedSubject = subject
			capturedMessage = message

			return nil
		},
	}

	// 設定とイベントサービスの準備
	cfg := &config.Config{
		Recording: struct {
			SaveDir string `toml:"save_dir" validate:"required"`
		}{
			SaveDir: tempDir,
		},
	}
	logger := zaptest.NewLogger(t)
	requestedService := eventutil.NewEventService[recording.RequestedEvent](
		mockBroker,
		logger,
		"recoto.recording.requested.v1",
	)

	// サービスの作成
	service := NewRecordingService(cfg, requestedService)

	// リクエストの作成
	req := &recordingv1.StartFromURLRequest{
		Url:      "http://example.com/stream",
		Title:    "Test_Title",
		Duration: durationpb.New(30 * time.Minute),
	}

	// リクエストの実行
	resp, err := service.StartFromURL(t.Context(), req)
	require.NoError(t, err)

	// 結果の検証
	assert.NotEmpty(t, resp.GetRecordingId())

	// イベントが正しく発行されたことを確認
	assert.Equal(t, "recoto.recording.requested.v1", capturedSubject)
	assert.NotEmpty(t, capturedMessage)
	// 出力ファイルパスが期待通りであることを確認
	assert.Contains(t, string(capturedMessage), "Test_Title.m4a")
}

func TestRecordingService_StartFromURL_FileExists(t *testing.T) {
	t.Parallel()

	// 一時ディレクトリの作成
	tempDir := t.TempDir()

	// テスト用のファイルを作成
	testFilePath := filepath.Join(tempDir, "Test_Title.m4a")
	err := os.WriteFile(testFilePath, []byte("test"), 0o600)
	require.NoError(t, err)

	// publishされたイベントをキャプチャするためのモック
	var (
		capturedSubject string
		capturedMessage []byte
	)

	mockBroker := &mockRecordingBroker{
		publishFunc: func(_ context.Context, subject string, message []byte) error {
			capturedSubject = subject
			capturedMessage = message

			return nil
		},
	}

	// 設定とイベントサービスの準備
	cfg := &config.Config{
		Recording: struct {
			SaveDir string `toml:"save_dir" validate:"required"`
		}{
			SaveDir: tempDir,
		},
	}
	logger := zaptest.NewLogger(t)
	requestedService := eventutil.NewEventService[recording.RequestedEvent](
		mockBroker,
		logger,
		"recoto.recording.requested.v1",
	)

	// サービスの作成
	service := NewRecordingService(cfg, requestedService)

	// リクエストの作成
	req := &recordingv1.StartFromURLRequest{
		Url:      "http://example.com/stream",
		Title:    "Test_Title",
		Duration: durationpb.New(30 * time.Minute),
	}

	// リクエストの実行
	resp, err := service.StartFromURL(t.Context(), req)
	require.NoError(t, err)

	// 結果の検証
	assert.NotEmpty(t, resp.GetRecordingId())

	// イベントが正しく発行されたことを確認
	assert.Equal(t, "recoto.recording.requested.v1", capturedSubject)
	assert.NotEmpty(t, capturedMessage)

	// ファイルパスにタイムスタンプが追加されていることを確認
	assert.NotContains(t, string(capturedMessage), testFilePath)
	assert.Contains(t, string(capturedMessage), "Test_Title_")
	assert.Contains(t, string(capturedMessage), ".m4a")
}

func TestRecordingService_StartFromURL_InvalidURL(t *testing.T) {
	t.Parallel()

	// 一時ディレクトリの作成
	tempDir := t.TempDir()

	// publishされたイベントをキャプチャするためのモック
	publishCalled := false
	mockBroker := &mockRecordingBroker{
		publishFunc: func(_ context.Context, _ string, _ []byte) error {
			publishCalled = true

			return nil
		},
	}

	// 設定とイベントサービスの準備
	cfg := &config.Config{
		Recording: struct {
			SaveDir string `toml:"save_dir" validate:"required"`
		}{
			SaveDir: tempDir,
		},
	}
	logger := zaptest.NewLogger(t)
	requestedService := eventutil.NewEventService[recording.RequestedEvent](
		mockBroker,
		logger,
		"recoto.recording.requested.v1",
	)

	// サービスの作成
	service := NewRecordingService(cfg, requestedService)

	// 空のURLでリクエストを作成（これは必ずエラーになる）
	req := &recordingv1.StartFromURLRequest{
		Url:      "",
		Title:    "Test Title",
		Duration: durationpb.New(30 * time.Minute),
	}

	// リクエストの実行
	resp, err := service.StartFromURL(t.Context(), req)
	require.Error(t, err)

	// 結果の検証
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to create")
	assert.False(t, publishCalled)
}

func TestRecordingService_StartFromURL_PublishError(t *testing.T) {
	t.Parallel()

	// 一時ディレクトリの作成
	tempDir := t.TempDir()

	// パブリッシュ時にエラーを返すモック
	publishErr := errors.New("publish error")
	mockBroker := &mockRecordingBroker{
		publishFunc: func(_ context.Context, _ string, _ []byte) error {
			return publishErr
		},
	}

	// 設定とイベントサービスの準備
	cfg := &config.Config{
		Recording: struct {
			SaveDir string `toml:"save_dir" validate:"required"`
		}{
			SaveDir: tempDir,
		},
	}
	logger := zaptest.NewLogger(t)
	requestedService := eventutil.NewEventService[recording.RequestedEvent](
		mockBroker,
		logger,
		"recoto.recording.requested.v1",
	)

	// サービスの作成
	service := NewRecordingService(cfg, requestedService)

	// リクエストの作成
	req := &recordingv1.StartFromURLRequest{
		Url:      "http://example.com/stream",
		Title:    "Test Title",
		Duration: durationpb.New(30 * time.Minute),
	}

	// リクエストの実行
	resp, err := service.StartFromURL(t.Context(), req)
	require.Error(t, err)

	// 結果の検証
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to publish requested event")
	assert.True(t, errors.Is(errors.Unwrap(errors.Unwrap(err)), publishErr))
}
