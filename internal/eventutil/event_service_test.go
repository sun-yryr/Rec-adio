package eventutil

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"

	"github.com/sun-yryr/recoto/internal/broker"
)

// testEvent はテスト用のイベント構造体.
type testEvent struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// mockBroker はブローカーのモック実装.
type mockBroker struct {
	publishFunc    func(ctx context.Context, subject string, message []byte) error
	subscribeFunc  func(ctx context.Context, subject string, handler func(message []byte)) (broker.UnsubscribeFunc, error)
	publishCalled  bool
	publishSubject string
	publishMessage []byte
	mu             sync.Mutex
}

func (m *mockBroker) Publish(ctx context.Context, subject string, message []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.publishCalled = true
	m.publishSubject = subject
	m.publishMessage = message

	if m.publishFunc != nil {
		return m.publishFunc(ctx, subject, message)
	}

	return nil
}

func (m *mockBroker) Subscribe(
	ctx context.Context,
	subject string,
	handler func(message []byte),
) (broker.UnsubscribeFunc, error) {
	if m.subscribeFunc != nil {
		return m.subscribeFunc(ctx, subject, handler)
	}

	return func() error { return nil }, nil
}

// Close はブローカーを閉じる.
func (m *mockBroker) Close() error {
	return nil
}

func (m *mockBroker) PublishCalled() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.publishCalled
}

func (m *mockBroker) PublishSubject() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.publishSubject
}

func (m *mockBroker) PublishMessage() []byte {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.publishMessage
}

func TestNewEventService(t *testing.T) {
	t.Parallel()

	// テスト用の入力
	testSubject := "test-subject"
	mockBroker := &mockBroker{}
	logger := zaptest.NewLogger(t)

	// テスト対象の関数を実行
	service := NewEventService[testEvent](mockBroker, logger, testSubject)

	// 結果の検証
	assert.NotNil(t, service)
	assert.Equal(t, mockBroker, service.broker)
	assert.Equal(t, logger, service.logger)
	assert.Equal(t, testSubject, service.subject)
}

func TestEventService_Publish(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		event         testEvent
		publishFunc   func(ctx context.Context, subject string, message []byte) error
		expectedError bool
	}{
		{
			name: "successful publish",
			event: testEvent{
				ID:      "1",
				Message: "test message",
			},
			publishFunc:   nil, // デフォルトの成功ケース
			expectedError: false,
		},
		{
			name: "publish error",
			event: testEvent{
				ID:      "2",
				Message: "will fail",
			},
			publishFunc: func(_ context.Context, _ string, _ []byte) error {
				return errors.New("failed to publish message")
			},
			expectedError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// モックとサービスの設定
			mockBroker := &mockBroker{
				publishFunc: testCase.publishFunc,
			}
			logger := zaptest.NewLogger(t)
			service := NewEventService[testEvent](mockBroker, logger, "test-subject")

			// テスト対象の関数を実行
			err := service.Publish(t.Context(), testCase.event)

			// 結果の検証
			if testCase.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.True(t, mockBroker.PublishCalled())
				assert.Equal(t, "test-subject", mockBroker.PublishSubject())
				// メッセージの内容を検証
				var decodedEvent testEvent

				require.NoError(t, json.Unmarshal(mockBroker.PublishMessage(), &decodedEvent))
				assert.Equal(t, testCase.event.ID, decodedEvent.ID)
				assert.Equal(t, testCase.event.Message, decodedEvent.Message)
			}
		})
	}
}

func TestEventService_Subscribe(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		subscribeFunc  func(ctx context.Context, subject string, handler func(message []byte)) (broker.UnsubscribeFunc, error)
		expectedError  bool
		messageToSend  []byte
		expectedID     string
		expectedMsg    string
		invalidMessage bool
	}{
		{
			name: "successful subscribe",
			subscribeFunc: func(_ context.Context, _ string, handler func(message []byte)) (broker.UnsubscribeFunc, error) {
				validEvent := testEvent{ID: "1", Message: "test message"}
				eventJSON, err := json.Marshal(validEvent)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal event: %w", err)
				}
				go func() {
					time.Sleep(100 * time.Millisecond) // 少し待って非同期処理
					handler(eventJSON)
				}()

				return func() error { return nil }, nil
			},
			expectedError: false,
			expectedID:    "1",
			expectedMsg:   "test message",
		},
		{
			name: "subscribe error",
			subscribeFunc: func(_ context.Context, _ string, _ func(message []byte)) (broker.UnsubscribeFunc, error) {
				return nil, errors.New("failed to subscribe to subject")
			},
			expectedError: true,
		},
		{
			name: "invalid message format",
			subscribeFunc: func(_ context.Context, _ string, handler func(message []byte)) (broker.UnsubscribeFunc, error) {
				go func() {
					time.Sleep(100 * time.Millisecond)
					handler([]byte("invalid json"))
				}()

				return func() error { return nil }, nil
			},
			expectedError:  false,
			invalidMessage: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// モックとサービスの設定
			mockBroker := &mockBroker{
				subscribeFunc: testCase.subscribeFunc,
			}

			// テスト用のロガー
			var logBuffer zaptest.Buffer

			logger := zaptest.NewLogger(
				t,
				zaptest.WrapOptions(zap.WrapCore(func(_ zapcore.Core) zapcore.Core {
					return zapcore.NewCore(
						zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
						&logBuffer,
						zapcore.DebugLevel,
					)
				})),
			)

			service := NewEventService[testEvent](mockBroker, logger, "test-subject")

			// イベント受信を検証するためのチャネル
			eventReceived := make(chan testEvent, 1)

			// テスト対象の関数を実行
			unsubscribe, err := service.Subscribe(
				t.Context(),
				func(_ context.Context, event *testEvent) {
					eventReceived <- *event
				},
			)

			// エラーケースの検証
			if testCase.expectedError {
				require.Error(t, err)
				assert.Nil(t, unsubscribe)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, unsubscribe)

			// 無効なメッセージの場合はログが出力されるだけで、ハンドラは呼ばれない
			if testCase.invalidMessage {
				time.Sleep(200 * time.Millisecond) // ログ出力を待つ
				assert.Contains(t, logBuffer.String(), "failed to unmarshal message")

				return
			}

			// 有効なメッセージの場合はハンドラが呼ばれる
			select {
			case event := <-eventReceived:
				assert.Equal(t, testCase.expectedID, event.ID)
				assert.Equal(t, testCase.expectedMsg, event.Message)
			case <-time.After(500 * time.Millisecond):
				t.Fatal("timeout waiting for event")
			}
		})
	}
}
