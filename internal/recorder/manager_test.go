package recorder

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/sun-yryr/recoto/internal/broker"
	"github.com/sun-yryr/recoto/internal/domain"
	"github.com/sun-yryr/recoto/internal/event/recording"
	"github.com/sun-yryr/recoto/internal/eventutil"
)

// mockRecorder はRecorderインターフェースのモック実装。
type mockRecorder struct {
	availableErr  error
	recErr        error
	supportSource []domain.SourceKind
	name          string
	recCalled     bool
	calledEvent   *recording.RequestedEvent
	recFunc       func(ctx context.Context, event *recording.RequestedEvent) error
}

func (m *mockRecorder) GetSupportSource() []domain.SourceKind {
	return m.supportSource
}

func (m *mockRecorder) Rec(ctx context.Context, event *recording.RequestedEvent) error {
	m.recCalled = true
	m.calledEvent = event

	if m.recFunc != nil {
		return m.recFunc(ctx, event)
	}

	return m.recErr
}

func (m *mockRecorder) GetName() string {
	return m.name
}

func (m *mockRecorder) CheckAvailable() error {
	return m.availableErr
}

var (
	errRecorderNotAvailable = errors.New("recorder not available error")
	errBrokerSubscribe      = errors.New("broker subscribe error")
)

type mockBroker struct {
	publishCalled    bool
	publishSubject   string
	publishMsg       []byte
	publishErr       error
	subscribeSubject string
	subscribeHandler func([]byte)
	subscribeErr     error
	unsubscribeErr   error
	closeErr         error
}

func (b *mockBroker) Publish(_ context.Context, subject string, msg []byte) error {
	b.publishCalled = true
	b.publishSubject = subject
	b.publishMsg = msg

	return b.publishErr
}

func (b *mockBroker) Subscribe(
	_ context.Context,
	subject string,
	handler func([]byte),
) (broker.UnsubscribeFunc, error) {
	b.subscribeSubject = subject
	b.subscribeHandler = handler

	return func() error { return b.unsubscribeErr }, b.subscribeErr
}

func (b *mockBroker) Close() error {
	return b.closeErr
}

func setupTestRecordingManager(t *testing.T) (*RecordingManager, *mockBroker) {
	t.Helper()

	logger := zaptest.NewLogger(t)
	mockBroker := &mockBroker{}

	requestedService := eventutil.NewEventService[recording.RequestedEvent](
		mockBroker,
		logger,
		"recording.requested",
	)
	startedService := eventutil.NewEventService[recording.StartedEvent](
		mockBroker,
		logger,
		"recording.started",
	)
	finishedService := eventutil.NewEventService[recording.FinishedEvent](
		mockBroker,
		logger,
		"recording.finished",
	)

	manager := NewRecordingManager(
		requestedService,
		startedService,
		finishedService,
		logger,
	)

	return manager, mockBroker
}

func TestNewRecordingManager(t *testing.T) {
	t.Parallel()

	manager, _ := setupTestRecordingManager(t)

	assert.NotNil(t, manager)
	assert.NotNil(t, manager.requestedService)
	assert.NotNil(t, manager.startedService)
	assert.NotNil(t, manager.finishedService)
	assert.Empty(t, manager.recorders)
	assert.Empty(t, manager.cancels)
}

func TestRecordingManager_AddRecorder(t *testing.T) {
	t.Parallel()

	manager, _ := setupTestRecordingManager(t)

	// 正常なレコーダーを追加
	recorder1 := &mockRecorder{
		name:          "Recorder1",
		supportSource: []domain.SourceKind{domain.SourceKindURL},
	}
	err := manager.AddRecorder(recorder1)
	require.NoError(t, err)

	assert.Len(t, manager.recorders, 1)
	assert.Equal(t, recorder1, manager.recorders[0].recorder)
	assert.True(t, manager.recorders[0].available)
	require.NoError(t, manager.recorders[0].error)

	// 利用できないレコーダーを追加
	recorder2 := &mockRecorder{
		name:          "Recorder2",
		supportSource: []domain.SourceKind{domain.SourceKindURL},
		availableErr:  errRecorderNotAvailable,
	}
	err = manager.AddRecorder(recorder2)
	require.Error(t, err)

	assert.Equal(t, errRecorderNotAvailable, err)
	assert.Len(t, manager.recorders, 2)
	assert.Equal(t, recorder2, manager.recorders[1].recorder)
	assert.False(t, manager.recorders[1].available)
	assert.Equal(t, errRecorderNotAvailable, manager.recorders[1].error)
}

func TestRecordingManager_Start(t *testing.T) {
	t.Parallel()

	manager, broker := setupTestRecordingManager(t)

	// ハンドラがnilにならないよう、デフォルト値を設定
	broker.subscribeHandler = func([]byte) {}

	// マネージャーの起動
	err := manager.Start(t.Context())
	require.NoError(t, err)

	// ブローカーのSubscribeが正しく呼ばれたことを確認
	assert.Equal(t, "recording.requested", broker.subscribeSubject)
	assert.NotNil(t, broker.subscribeHandler)

	// エラーケースのテスト
	errorBroker := &mockBroker{
		subscribeErr: errBrokerSubscribe,
	}
	logger := zaptest.NewLogger(t)

	requestedService := eventutil.NewEventService[recording.RequestedEvent](
		errorBroker,
		logger,
		"recording.requested",
	)
	startedService := eventutil.NewEventService[recording.StartedEvent](
		errorBroker,
		logger,
		"recording.started",
	)
	finishedService := eventutil.NewEventService[recording.FinishedEvent](
		errorBroker,
		logger,
		"recording.finished",
	)

	manager2 := NewRecordingManager(
		requestedService,
		startedService,
		finishedService,
		logger,
	)

	err = manager2.Start(t.Context())
	require.Error(t, err)

	assert.Contains(t, err.Error(), "subscribe error")
}

func TestRecordingManager_CancelRecording(t *testing.T) {
	t.Parallel()

	// テスト用のマネージャー
	manager, _ := setupTestRecordingManager(t)

	// レコーディングIDとダミーのキャンセル関数
	recordingID := "test-recording-id"
	cancelCalled := false
	cancelFunc := func() {
		cancelCalled = true
	}

	// キャンセル関数を登録
	manager.mu.Lock()
	manager.cancels[recordingID] = cancelFunc
	manager.mu.Unlock()

	// 存在するIDをキャンセル
	cancelled := manager.CancelRecording(recordingID)
	assert.True(t, cancelled, "should return true when cancelling an existing recording")
	assert.True(t, cancelCalled, "cancel function should be called")

	// キャンセル後にmapからエントリが削除されていることを確認
	manager.mu.Lock()
	_, exists := manager.cancels[recordingID]
	manager.mu.Unlock()
	assert.False(t, exists, "recording should be removed from cancels map")

	// 存在しないIDをキャンセル
	cancelled = manager.CancelRecording("non-existent-id")
	assert.False(t, cancelled, "should return false when cancelling a non-existent recording")
}

func TestRecordingManager_ProcessRecording_Cancellation(t *testing.T) {
	t.Parallel()

	// テスト用のマネージャー
	manager, _ := setupTestRecordingManager(t)

	// コンテキストを作成
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	// テスト用のイベント
	source, err := domain.NewURLSource("http://example.com/stream")
	require.NoError(t, err)

	event := &recording.RequestedEvent{
		RecordingID: "test-id",
		Source:      source,
		Output:      "test-output.m4a",
		Duration:    30 * time.Second,
	}

	// このRecは録音中にコンテキストのキャンセルを検出するモック
	recCalled := false
	testMutex := sync.Mutex{}

	recorder := &mockRecorder{
		name:          "TestRecorder",
		supportSource: []domain.SourceKind{domain.SourceKindURL},
		recFunc: func(ctx context.Context, _ *recording.RequestedEvent) error {
			testMutex.Lock()
			recCalled = true
			testMutex.Unlock()
			// 録音中にキャンセルを待つ
			<-ctx.Done()

			return context.Canceled
		},
	}

	// キャンセルを登録
	manager.cancels[event.RecordingID] = cancel

	// 新しいgoroutineで録音プロセスを実行
	processDone := make(chan struct{})
	go func() {
		manager.processRecording(ctx, event, recorder)
		close(processDone)
	}()

	// 録音関数が呼ばれるのを待つ - 最大100ミリ秒待機
	for i := range 10 {
		testMutex.Lock()
		called := recCalled
		testMutex.Unlock()

		if called {
			break
		}

		if i < 9 {
			time.Sleep(10 * time.Millisecond)
		}
	}

	// キャンセルする
	cancel()

	// プロセスが完了するのを待つ
	select {
	case <-processDone:
		// プロセスが完了した
	case <-time.After(1 * time.Second):
		t.Fatal("Process did not complete within timeout")
	}

	// レコーダーが呼ばれたことを確認
	testMutex.Lock()
	assert.True(t, recCalled, "Recorder.Rec should have been called")
	testMutex.Unlock()
}
