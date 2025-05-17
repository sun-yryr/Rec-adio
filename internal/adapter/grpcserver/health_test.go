package grpcserver

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/sun-yryr/recoto/internal/broker"
	"github.com/sun-yryr/recoto/internal/logger"
	pb "github.com/sun-yryr/recoto/pkg/api/recoto/health/v1"
)

// mockHealthChecker はヘルスチェッカーのモック
type mockHealthChecker struct {
	name      string
	checkFunc func(ctx context.Context) error
}

func (m *mockHealthChecker) GetName() string {
	return m.name
}

func (m *mockHealthChecker) Check(ctx context.Context) error {
	if m.checkFunc != nil {
		return m.checkFunc(ctx)
	}
	return nil
}

// mockBroker はブローカーのモック
type mockBroker struct{}

func (m *mockBroker) Publish(_ context.Context, _ string, _ []byte) error {
	return nil
}

func (m *mockBroker) Subscribe(
	_ context.Context,
	_ string,
	_ func(message []byte),
) (broker.UnsubscribeFunc, error) {
	return func() error { return nil }, nil
}

func (m *mockBroker) Close() error {
	return nil
}

func TestNewHealthService(t *testing.T) {
	t.Parallel()

	// モックの準備
	mockBroker := &mockBroker{}
	checker1 := &mockHealthChecker{name: "checker1"}
	checker2 := &mockHealthChecker{name: "checker2"}

	// サービスの作成
	service := NewHealthService(mockBroker, checker1, checker2)

	// 結果の検証
	assert.NotNil(t, service)
	assert.Equal(t, mockBroker, service.broker)
	assert.Len(t, service.healthCheckers, 2)
	assert.Equal(t, checker1, service.healthCheckers[0])
	assert.Equal(t, checker2, service.healthCheckers[1])
}

func TestHealthService_Check_AllOK(t *testing.T) {
	t.Parallel()

	// テスト用のロガーとコンテキストを準備
	testLogger := zaptest.NewLogger(t)
	ctx := logger.WithLogger(context.Background(), testLogger)

	// 常に成功するチェッカーを作成
	checker1 := &mockHealthChecker{name: "checker1"}
	checker2 := &mockHealthChecker{name: "checker2"}

	// サービスの作成
	service := NewHealthService(&mockBroker{}, checker1, checker2)

	// ヘルスチェックの実行
	resp, err := service.Check(ctx, &pb.CheckRequest{})

	// 結果の検証
	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Status)
	assert.Len(t, resp.Results, 2)

	// すべてのチェッカーが成功していることを確認
	for _, result := range resp.Results {
		assert.True(t, result.Ok)
		assert.Empty(t, result.Error)
	}
}

func TestHealthService_Check_PartialFailure(t *testing.T) {
	t.Parallel()

	// テスト用のロガーとコンテキストを準備
	testLogger := zaptest.NewLogger(t)
	ctx := logger.WithLogger(context.Background(), testLogger)

	// 1つは成功、1つは失敗するチェッカーを作成
	checker1 := &mockHealthChecker{name: "checker1"}
	testError := errors.New("test error")
	checker2 := &mockHealthChecker{
		name: "checker2",
		checkFunc: func(ctx context.Context) error {
			return testError
		},
	}

	// サービスの作成
	service := NewHealthService(&mockBroker{}, checker1, checker2)

	// ヘルスチェックの実行
	resp, err := service.Check(ctx, &pb.CheckRequest{})

	// 結果の検証
	require.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Len(t, resp.Results, 2)

	// 結果の詳細を確認
	var successCount, failureCount int
	for _, result := range resp.Results {
		if result.Ok {
			successCount++
		} else {
			failureCount++
			assert.Contains(t, result.Error, testError.Error())
		}
	}

	assert.Equal(t, 1, successCount)
	assert.Equal(t, 1, failureCount)
}

func TestHealthService_Check_AllFail(t *testing.T) {
	t.Parallel()

	// テスト用のロガーとコンテキストを準備
	testLogger := zaptest.NewLogger(t)
	ctx := logger.WithLogger(context.Background(), testLogger)

	// すべて失敗するチェッカーを作成
	testError1 := errors.New("test error 1")
	checker1 := &mockHealthChecker{
		name: "checker1",
		checkFunc: func(ctx context.Context) error {
			return testError1
		},
	}
	testError2 := errors.New("test error 2")
	checker2 := &mockHealthChecker{
		name: "checker2",
		checkFunc: func(ctx context.Context) error {
			return testError2
		},
	}

	// サービスの作成
	service := NewHealthService(&mockBroker{}, checker1, checker2)

	// ヘルスチェックの実行
	resp, err := service.Check(ctx, &pb.CheckRequest{})

	// 結果の検証
	require.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Len(t, resp.Results, 2)

	// すべてのチェッカーが失敗していることを確認
	var foundError1, foundError2 bool
	for _, result := range resp.Results {
		assert.False(t, result.Ok)
		if result.Error == testError1.Error() {
			foundError1 = true
		}
		if result.Error == testError2.Error() {
			foundError2 = true
		}
	}
	assert.True(t, foundError1, "Expected to find error with message '%s'", testError1.Error())
	assert.True(t, foundError2, "Expected to find error with message '%s'", testError2.Error())
}

func TestHealthService_Check_NoCheckers(t *testing.T) {
	t.Parallel()

	// テスト用のロガーとコンテキストを準備
	testLogger := zaptest.NewLogger(t)
	ctx := logger.WithLogger(context.Background(), testLogger)

	// チェッカーなしでサービスを作成
	service := NewHealthService(&mockBroker{})

	// ヘルスチェックの実行
	resp, err := service.Check(ctx, &pb.CheckRequest{})

	// 結果の検証
	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Status)
	assert.Empty(t, resp.Results)
}