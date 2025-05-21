//nolint:paralleltest // 内臓Brokerを使っているので並列にしない
package nats

import (
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestBroker_PublishSubscribe(t *testing.T) {
	if testing.Short() {
		t.Skip("短いテストモードでは埋め込みNATSサーバーのテストをスキップ")
	}

	t.Run("PublishSubscribe", func(tt *testing.T) {
		logger := zaptest.NewLogger(tt)
		broker, err := NewEmbeddedBroker(logger)
		require.NoError(tt, err)

		// require.Eventuallyで100msごとに5秒間起動を待つ
		require.Eventually(tt, func() bool {
			return broker.nc != nil && broker.nc.Status() == nats.CONNECTED
		}, 5*time.Second, 100*time.Millisecond, "timed out waiting for NATS connection")

		defer func() {
			err := broker.Close()
			assert.NoError(tt, err)
		}()

		subject := "test.publish.subscribe"
		message := []byte("hello world")

		var wg sync.WaitGroup

		wg.Add(1)

		var receivedMessage []byte

		unsubscribe, err := broker.Subscribe(tt.Context(), subject, func(msg []byte) {
			receivedMessage = msg

			wg.Done()
		})
		require.NoError(tt, err)

		defer func() {
			err := unsubscribe()
			assert.NoError(tt, err)
		}()

		err = broker.Publish(tt.Context(), subject, message)
		require.NoError(tt, err)

		// メッセージが受信されるのを待つ
		wg.Wait()

		assert.Equal(tt, message, receivedMessage)
	})

	t.Run("Close", func(tt *testing.T) {
		logger := zaptest.NewLogger(tt)
		broker, err := NewEmbeddedBroker(logger)
		require.NoError(tt, err)

		// Wait for NATS server to be ready using require.Eventually
		require.Eventually(tt, func() bool {
			return broker.nc != nil && broker.nc.Status() == nats.CONNECTED
		}, 5*time.Second, 100*time.Millisecond, "timed out waiting for NATS connection")

		// Closeが正常に完了することを確認
		err = broker.Close()
		require.NoError(tt, err)

		// 閉じた後のコネクションのステータスを確認
		assert.Equal(tt, nats.CLOSED, broker.nc.Status())
	})

	t.Run("Error when closed broker", func(tt *testing.T) {
		logger := zaptest.NewLogger(tt)
		broker, err := NewEmbeddedBroker(logger)
		require.NoError(tt, err)

		// Wait for NATS server to be ready using require.Eventually
		require.Eventually(tt, func() bool {
			return broker.nc != nil && broker.nc.Status() == nats.CONNECTED
		}, 5*time.Second, 100*time.Millisecond, "timed out waiting for NATS connection")

		// ブローカーを閉じる
		err = broker.Close()
		require.NoError(tt, err)

		// 閉じたブローカーでのパブリッシュは失敗するはず
		err = broker.Publish(tt.Context(), "test.subject", []byte("test message"))
		require.Error(tt, err)

		// 閉じたブローカーでのサブスクライブは失敗するはず
		_, err = broker.Subscribe(tt.Context(), "test.subject", func(_ []byte) {})
		assert.Error(tt, err)
	})
}
