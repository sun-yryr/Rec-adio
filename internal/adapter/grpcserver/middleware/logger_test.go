package middleware

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"

	"github.com/sun-yryr/recoto/internal/logger"
)

func TestNewLoggerInterceptor(t *testing.T) {
	t.Parallel()

	// テスト用のロガーを作成
	testLogger := zaptest.NewLogger(t)

	// インターセプタを作成
	interceptor := NewLoggerInterceptor(testLogger)

	// テスト用のリクエストとハンドラ
	testReq := "test-request"
	testHandler := func(ctx context.Context, _ interface{}) (interface{}, error) {
		log := logger.FromContext(ctx)
		assert.NotEmpty(t, log.Core()) // NewNopで生成されていないことを検証

		return "test-response", nil
	}

	resp, err := interceptor(t.Context(), testReq, &grpc.UnaryServerInfo{}, testHandler)
	require.NoError(t, err)

	assert.Equal(t, "test-response", resp)
}
