package broker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockBroker struct {
	publishFunc    func(ctx context.Context, subject string, message interface{}) error
	subscribeFunc  func(ctx context.Context, subject string, handler func(context.Context, []byte)) (UnsubscribeFunc, error)
	closeFunc      func() error
	publishCalled  bool
	publishSubject string
	publishMessage interface{}
}

func (m *mockBroker) Publish(ctx context.Context, subject string, message interface{}) error {
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
	handler func(context.Context, []byte),
) (UnsubscribeFunc, error) {
	if m.subscribeFunc != nil {
		return m.subscribeFunc(ctx, subject, handler)
	}

	return func() error { return nil }, nil
}

func (m *mockBroker) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}

	return nil
}

func TestNewHealthCheck(t *testing.T) {
	t.Parallel()

	mockBroker := &mockBroker{}
	hc := NewHealthCheck(mockBroker)

	assert.NotNil(t, hc)
	assert.Equal(t, mockBroker, hc.Broker)
}

func TestHealthCheck_GetName(t *testing.T) {
	t.Parallel()

	mockBroker := &mockBroker{}
	hc := NewHealthCheck(mockBroker)

	assert.Equal(t, "BrokerHealthCheck", hc.GetName())
}

func TestHealthCheck_Check_Success(t *testing.T) {
	t.Parallel()

	// 成功するケース: Subscribeが成功し、パブリッシュされたメッセージに応答する
	mockBroker := &mockBroker{
		subscribeFunc: func(ctx context.Context, _ string, handler func(context.Context, []byte)) (UnsubscribeFunc, error) {
			// パブリッシュされると同時にハンドラをトリガーするモック
			go func(ctx context.Context) {
				time.Sleep(50 * time.Millisecond) // 少し遅延を入れる
				handler(ctx, []byte("pong"))
			}(ctx)

			return func() error { return nil }, nil
		},
	}

	hc := NewHealthCheck(mockBroker)
	err := hc.Check(t.Context())
	require.NoError(t, err)

	assert.True(t, mockBroker.publishCalled)
	assert.Equal(t, "recoto.health.check.v1", mockBroker.publishSubject)
	assert.Equal(t, "ping", mockBroker.publishMessage)
}

func TestHealthCheck_Check_SubscribeError(t *testing.T) {
	t.Parallel()

	// Subscribeがエラーを返すケース
	errSubscribe := assert.AnError
	mockBroker := &mockBroker{
		subscribeFunc: func(_ context.Context, _ string, _ func(context.Context, []byte)) (UnsubscribeFunc, error) {
			return nil, errSubscribe
		},
	}

	hc := NewHealthCheck(mockBroker)
	err := hc.Check(t.Context())
	require.Error(t, err)

	require.ErrorIs(t, err, errSubscribe)
	assert.False(t, mockBroker.publishCalled)
}

func TestHealthCheck_Check_PublishError(t *testing.T) {
	t.Parallel()

	// Publishがエラーを返すケース
	errPublish := assert.AnError
	mockBroker := &mockBroker{
		subscribeFunc: func(_ context.Context, _ string, _ func(context.Context, []byte)) (UnsubscribeFunc, error) {
			return func() error { return nil }, nil
		},
		publishFunc: func(_ context.Context, _ string, _ interface{}) error {
			return errPublish
		},
	}

	hc := NewHealthCheck(mockBroker)
	err := hc.Check(t.Context())
	require.Error(t, err)

	assert.ErrorIs(t, err, errPublish)
}

func TestHealthCheck_Check_Timeout(t *testing.T) {
	t.Parallel()

	// タイムアウトが発生するケース
	mockBroker := &mockBroker{
		subscribeFunc: func(_ context.Context, _ string, _ func(context.Context, []byte)) (UnsubscribeFunc, error) {
			// ハンドラを呼び出さない
			return func() error { return nil }, nil
		},
	}

	hc := NewHealthCheck(mockBroker)

	// タイムアウトを短く設定したコンテキストを使用
	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()

	err := hc.Check(ctx)
	require.Error(t, err)

	assert.True(t, mockBroker.publishCalled)
	assert.Contains(t, err.Error(), "health check timeout")
}

func TestHealthCheck_Check_UnsubscribeError(t *testing.T) {
	t.Parallel()

	// Unsubscribeがエラーを返すケース（ただしヘルスチェック自体は成功する）
	unsubscribeErr := assert.AnError
	unsubscribeCalled := false
	handlerCalled := false

	mockBroker := &mockBroker{
		subscribeFunc: func(ctx context.Context, _ string, handler func(context.Context, []byte)) (UnsubscribeFunc, error) {
			// パブリッシュされると同時にハンドラをトリガーするモック
			go func(ctx context.Context) {
				time.Sleep(50 * time.Millisecond)
				handlerCalled = true
				handler(ctx, []byte("pong"))
			}(ctx)

			return func() error {
				unsubscribeCalled = true

				return unsubscribeErr
			}, nil
		},
	}

	hc := NewHealthCheck(mockBroker)
	err := hc.Check(t.Context())
	require.NoError(t, err)

	assert.True(t, handlerCalled)
	assert.True(t, unsubscribeCalled)
}
