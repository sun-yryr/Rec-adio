package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWithLoggerAndFromContext(t *testing.T) {
	t.Parallel()

	// テスト用のロガーを作成
	testLogger, _ := zap.NewProduction()
	defer func() { _ = testLogger.Sync() }()

	// コンテキストにロガーを追加
	ctx := t.Context()
	ctxWithLogger := WithLogger(ctx, testLogger)

	// コンテキストからロガーを取得
	logger := FromContext(ctxWithLogger)

	// 取得したロガーが元のロガーと同じインスタンスであることを確認
	assert.Same(t, testLogger, logger)
}

func TestFromContextWithNoLogger(t *testing.T) {
	t.Parallel()

	// ロガーを持たないコンテキスト
	ctx := t.Context()

	// コンテキストからロガーを取得
	logger := FromContext(ctx)

	// 代替のNopロガーが返されることを確認
	assert.NotNil(t, logger)
	// 実際にログ出力してもエラーが発生しないことを確認
	assert.NotPanics(t, func() {
		logger.Info("test message")
	})
}
