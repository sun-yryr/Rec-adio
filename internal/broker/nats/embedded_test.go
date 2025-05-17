package nats

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartEmbeddedServer(t *testing.T) {
	if testing.Short() {
		t.Skip("短いテストモードでは埋め込みNATSサーバーのテストをスキップ")
	}

	// NATSサーバーの起動
	server, err := startEmbeddedServer()
	require.NoError(t, err)
	defer server.Shutdown()

	// サーバーが正常に稼働していることを確認
	assert.True(t, server.Running())
	assert.NotEmpty(t, server.ClientURL())

	// サーバーが空でないポートで待ち受けていることを確認
	assert.NotEmpty(t, server.ClientURL())
}

func TestStartEmbeddedServer_ReadyForConnections(t *testing.T) {
	if testing.Short() {
		t.Skip("短いテストモードでは埋め込みNATSサーバーのテストをスキップ")
	}

	// NATSサーバーの起動
	server, err := startEmbeddedServer()
	require.NoError(t, err)
	defer server.Shutdown()

	// 接続準備完了していることを確認
	isReady := server.ReadyForConnections(1 * time.Second)
	assert.True(t, isReady)
}