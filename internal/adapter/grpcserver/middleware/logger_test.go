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
	require.NotNil(t, interceptor)

	// テスト用のリクエストとハンドラ
	testReq := "test-request"
	var handlerCtx context.Context
	testHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		handlerCtx = ctx
		return "test-response", nil
	}

	// インターセプタを実行
	resp, err := interceptor(context.Background(), testReq, &grpc.UnaryServerInfo{}, testHandler)

	// 結果の検証
	require.NoError(t, err)
	assert.Equal(t, "test-response", resp)

	// コンテキストにロガーが追加されていることを確認
	ctxLogger := logger.FromContext(handlerCtx)
	assert.NotNil(t, ctxLogger)
}