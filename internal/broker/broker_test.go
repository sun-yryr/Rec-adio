package broker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBrokerInterface(t *testing.T) {
	t.Parallel()

	// これは実質的なテストではなく、インターフェースの実装が型安全であることを確認するためのテスト
	var _ Broker = (*mockBroker)(nil)

	// mockBrokerがBrokerインターフェースを正しく実装していることを確認
	broker := &mockBroker{}

	// 各メソッドがエラーなく呼び出せることを確認
	err := broker.Publish(t.Context(), "test.subject", []byte("test message"))
	require.NoError(t, err)

	unsubscribe, err := broker.Subscribe(
		t.Context(),
		"test.subject",
		func(_ []byte) {},
	)
	require.NoError(t, err)
	assert.NotNil(t, unsubscribe)

	err = unsubscribe()
	require.NoError(t, err)

	err = broker.Close()
	require.NoError(t, err)
}
