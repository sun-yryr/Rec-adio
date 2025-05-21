package nats

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartEmbeddedServer(t *testing.T) {
	t.Parallel()

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
	assert.True(t, server.ReadyForConnections(1*time.Second))
}
