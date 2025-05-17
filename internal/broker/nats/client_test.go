package nats

import (
	"context"
	"sync"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// TestBroker_Interface はBrokerが必要なインターフェースを満たすかテスト
func TestBroker_Interface(t *testing.T) {
	// コンパイル時に満たされているか確認するだけなので実行はされない
	// 型アサーションをここで確認する
	var _ interface {
		Publish(ctx context.Context, subject string, message []byte) error
		Close() error
	} = (*Broker)(nil)

	// Subscribe型の確認は別に行う
	t.SkipNow() // 実際のテストコードは実行しない
}

// TestBroker_PublishSubscribe はPublishとSubscribeが連携して動作するかテスト
func TestBroker_PublishSubscribe(t *testing.T) {
	if testing.Short() {
		t.Skip("短いテストモードでは埋め込みNATSサーバーのテストをスキップ")
	}

	logger := zaptest.NewLogger(t)
	broker, err := NewEmbeddedBroker(logger)
	require.NoError(t, err)
	defer func() {
		err := broker.Close()
		assert.NoError(t, err)
	}()

	subject := "test.publish.subscribe"
	message := []byte("hello world")
	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(1)

	var receivedMessage []byte
	unsubscribe, err := broker.Subscribe(ctx, subject, func(msg []byte) {
		receivedMessage = msg
		wg.Done()
	})
	require.NoError(t, err)
	defer func() {
		err := unsubscribe()
		assert.NoError(t, err)
	}()

	err = broker.Publish(ctx, subject, message)
	require.NoError(t, err)

	// メッセージが受信されるのを待つ
	wg.Wait()

	assert.Equal(t, message, receivedMessage)
}

// TestBroker_Close はCloseが正しく動作するかテスト
func TestBroker_Close(t *testing.T) {
	if testing.Short() {
		t.Skip("短いテストモードでは埋め込みNATSサーバーのテストをスキップ")
	}

	logger := zaptest.NewLogger(t)
	broker, err := NewEmbeddedBroker(logger)
	require.NoError(t, err)

	// Closeが正常に完了することを確認
	err = broker.Close()
	assert.NoError(t, err)

	// 閉じた後のコネクションのステータスを確認
	assert.Equal(t, nats.CLOSED, broker.nc.Status())
}

// TestBroker_Errors はエラーケースをテスト
func TestBroker_Errors(t *testing.T) {
	if testing.Short() {
		t.Skip("短いテストモードでは埋め込みNATSサーバーのテストをスキップ")
	}

	t.Run("正常に開始したブローカーに対する操作", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		broker, err := NewEmbeddedBroker(logger)
		require.NoError(t, err)
		defer broker.Close()

		// 有効なブローカーでのサブスクライブとパブリッシュは成功するはず
		unsubscribe, err := broker.Subscribe(context.Background(), "test.subject", func(message []byte) {})
		assert.NoError(t, err)
		assert.NotNil(t, unsubscribe)

		err = broker.Publish(context.Background(), "test.subject", []byte("test message"))
		assert.NoError(t, err)
	})

	t.Run("閉じたブローカーに対する操作", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		broker, err := NewEmbeddedBroker(logger)
		require.NoError(t, err)

		// ブローカーを閉じる
		err = broker.Close()
		require.NoError(t, err)

		// 閉じたブローカーでのパブリッシュは失敗するはず
		err = broker.Publish(context.Background(), "test.subject", []byte("test message"))
		assert.Error(t, err)

		// 閉じたブローカーでのサブスクライブは失敗するはず
		_, err = broker.Subscribe(context.Background(), "test.subject", func(message []byte) {})
		assert.Error(t, err)
	})
}