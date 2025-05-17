package broker

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBrokerInterface(t *testing.T) {
	t.Parallel()

	// これは実質的なテストではなく、インターフェースの実装が型安全であることを確認するためのテスト
	var _ Broker = (*mockBroker)(nil)

	// mockBrokerがBrokerインターフェースを正しく実装していることを確認
	broker := &mockBroker{}
	
	// 各メソッドがエラーなく呼び出せることを確認
	err := broker.Publish(context.Background(), "test.subject", []byte("test message"))
	assert.NoError(t, err)

	unsubscribe, err := broker.Subscribe(context.Background(), "test.subject", func(message []byte) {})
	assert.NoError(t, err)
	assert.NotNil(t, unsubscribe)

	err = unsubscribe()
	assert.NoError(t, err)

	err = broker.Close()
	assert.NoError(t, err)
}