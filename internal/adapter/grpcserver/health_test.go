package grpcserver

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/sun-yryr/recoto/internal/broker"
	"github.com/sun-yryr/recoto/internal/logger"
	pb "github.com/sun-yryr/recoto/pkg/api/recoto/health/v1"
)

// mockHealthChecker はヘルスチェッカーのモック。
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

var (
	errTest  = errors.New("test error")
	errTest1 = errors.New("test error 1")
	errTest2 = errors.New("test error 2")
)

// mockBroker はブローカーのモック。
type mockBroker struct {
	published  bool
	subscribed bool
}

func (m *mockBroker) Publish(_ context.Context, _ string, _ []byte) error {
	m.published = true

	return nil
}

func (m *mockBroker) Subscribe(
	_ context.Context,
	_ string,
	_ func(message []byte),
) (broker.UnsubscribeFunc, error) {
	m.subscribed = true

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
	ctx := logger.WithLogger(t.Context(), testLogger)

	// 常に成功するチェッカーを作成
	checker1 := &mockHealthChecker{name: "checker1"}
	checker2 := &mockHealthChecker{name: "checker2"}

	// サービスの作成
	service := NewHealthService(&mockBroker{}, checker1, checker2)

	// ヘルスチェックの実行
	resp, err := service.Check(ctx, &pb.CheckRequest{})
	require.NoError(t, err)

	// 結果の検証
	assert.Equal(t, "ok", resp.GetStatus())
	assert.Len(t, resp.GetResults(), 2)

	// すべてのチェッカーが成功していることを確認
	for _, result := range resp.GetResults() {
		assert.True(t, result.GetOk())
		assert.Empty(t, result.GetError())
	}
}

func TestHealthService_Check_PartialFailure(t *testing.T) {
	t.Parallel()

	// テスト用のロガーとコンテキストを準備
	testLogger := zaptest.NewLogger(t)
	ctx := logger.WithLogger(t.Context(), testLogger)

	// 1つは成功、1つは失敗するチェッカーを作成
	checker1 := &mockHealthChecker{name: "checker1"}
	checker2 := &mockHealthChecker{
		name: "checker2",
		checkFunc: func(_ context.Context) error {
			return errTest
		},
	}

	// サービスの作成
	service := NewHealthService(&mockBroker{}, checker1, checker2)

	// ヘルスチェックの実行
	resp, err := service.Check(ctx, &pb.CheckRequest{})
	require.NoError(t, err)

	// 結果の検証
	assert.Equal(t, "error", resp.GetStatus())
	assert.Len(t, resp.GetResults(), 2)

	// 結果の詳細を確認
	var successCount, failureCount int

	for _, result := range resp.GetResults() {
		if result.GetOk() {
			successCount++
		} else {
			failureCount++

			assert.Contains(t, result.GetError(), errTest.Error())
		}
	}

	assert.Equal(t, 1, successCount)
	assert.Equal(t, 1, failureCount)
}

func TestHealthService_Check_AllFail(t *testing.T) {
	t.Parallel()

	// テスト用のロガーとコンテキストを準備
	testLogger := zaptest.NewLogger(t)
	ctx := logger.WithLogger(t.Context(), testLogger)

	// すべて失敗するチェッカーを作成
	checker1 := &mockHealthChecker{
		name: "checker1",
		checkFunc: func(_ context.Context) error {
			return errTest1
		},
	}
	checker2 := &mockHealthChecker{
		name: "checker2",
		checkFunc: func(_ context.Context) error {
			return errTest2
		},
	}

	// サービスの作成
	service := NewHealthService(&mockBroker{}, checker1, checker2)

	// ヘルスチェックの実行
	resp, err := service.Check(ctx, &pb.CheckRequest{})

	// 結果の検証
	require.NoError(t, err)
	assert.Equal(t, "error", resp.GetStatus())
	assert.Len(t, resp.GetResults(), 2)

	// すべてのチェッカーが失敗していることを確認
	var foundError1, foundError2 bool

	for _, result := range resp.GetResults() {
		assert.False(t, result.GetOk())

		if result.GetError() == errTest1.Error() {
			foundError1 = true
		}

		if result.GetError() == errTest2.Error() {
			foundError2 = true
		}
	}

	assert.True(t, foundError1, "Expected to find error with message '%s'", errTest1.Error())
	assert.True(t, foundError2, "Expected to find error with message '%s'", errTest2.Error())
}

func TestHealthService_Check_NoCheckers(t *testing.T) {
	t.Parallel()

	// テスト用のロガーとコンテキストを準備
	testLogger := zaptest.NewLogger(t)
	ctx := logger.WithLogger(t.Context(), testLogger)

	// チェッカーなしでサービスを作成
	service := NewHealthService(&mockBroker{})

	// ヘルスチェックの実行
	resp, err := service.Check(ctx, &pb.CheckRequest{})

	// 結果の検証
	require.NoError(t, err)
	assert.Equal(t, "ok", resp.GetStatus())
	assert.Empty(t, resp.GetResults())
}
