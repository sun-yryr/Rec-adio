package recorder

import (
	"context"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/sun-yryr/recoto/internal/domain"
	"github.com/sun-yryr/recoto/internal/event/recording"
)

func TestNewURLRecorder(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	recorder := NewURLRecorder(logger)

	assert.NotNil(t, recorder)
	assert.Equal(t, logger, recorder.logger)
}

func TestURLRecorder_GetSupportSource(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	recorder := NewURLRecorder(logger)

	sources := recorder.GetSupportSource()
	
	assert.Equal(t, []domain.SourceKind{domain.SourceKindURL}, sources)
}

func TestURLRecorder_GetName(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	recorder := NewURLRecorder(logger)

	assert.Equal(t, "URLRecorder", recorder.GetName())
}

func TestURLRecorder_CheckAvailable(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	recorder := NewURLRecorder(logger)

	// ffmpegが存在する場合はエラーなし
	if _, err := exec.LookPath("ffmpeg"); err == nil {
		assert.NoError(t, recorder.CheckAvailable())
	} else {
		// ffmpegが存在しない場合はエラー
		err := recorder.CheckAvailable()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ffmpeg is not installed")
	}
}

func TestURLRecorder_Rec_InvalidURL(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	recorder := NewURLRecorder(logger)

	// 無効なURLでテスト（完全に無効な形式）
	source := &domain.Source{
		Kind: domain.SourceKindURL,
		ID:   "://invalid",
	}

	event := &recording.RequestedEvent{
		RecordingID: "test-id",
		Source:      source,
		Output:      "test-output.m4a",
		Duration:    5 * time.Second,
	}

	err := recorder.Rec(context.Background(), event)
	require.Error(t, err, "Expected error for invalid URL")
	assert.Contains(t, err.Error(), "invalid URL")
}

func TestURLRecorder_Rec_TooLongDuration(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	recorder := NewURLRecorder(logger)

	// 非常に長い録音時間でテスト
	source := &domain.Source{
		Kind: domain.SourceKindURL,
		ID:   "http://example.com/valid.mp3",
	}

	event := &recording.RequestedEvent{
		RecordingID: "test-id",
		Source:      source,
		Output:      "test-output.m4a",
		Duration:    time.Duration(math.MaxInt64) - 5*time.Second, // 長すぎる時間
	}

	err := recorder.Rec(context.Background(), event)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "recording duration is too long")
}

func TestURLRecorder_Rec_Canceled(t *testing.T) {
	t.Parallel()

	// このテストはFFmpegがインストールされている場合のみ実行する
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg is not installed, skipping test")
	}

	logger := zaptest.NewLogger(t)
	recorder := NewURLRecorder(logger)

	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "output.m4a")

	// 存在するが録音に失敗するURLでテスト（キャンセルするため）
	source := &domain.Source{
		Kind: domain.SourceKindURL,
		ID:   "http://example.com/nonexistent.mp3", // 存在しないURLだがパースは成功する
	}

	event := &recording.RequestedEvent{
		RecordingID: "test-id",
		Source:      source,
		Output:      outputFile,
		Duration:    30 * time.Second,
	}

	// ちょっと待ってからキャンセルするコンテキスト
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	// キャンセルされるのでエラーになるはず
	err := recorder.Rec(ctx, event)
	assert.Error(t, err)
}

// この統合テストは、実際のURLからの録音を試みる
// 環境変数TEST_INTEGRATION=trueの場合のみ実行する
func TestURLRecorder_Rec_Integration(t *testing.T) {
	if os.Getenv("TEST_INTEGRATION") != "true" {
		t.Skip("skipping integration test; set TEST_INTEGRATION=true to run")
	}

	// このテストはFFmpegがインストールされている場合のみ実行する
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg is not installed, skipping test")
	}

	logger := zaptest.NewLogger(t)
	recorder := NewURLRecorder(logger)

	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "output.m4a")

	// 実際に録音可能なURLを指定（このURLは適宜変更する必要がある）
	// テスト用のストリーミングサーバー: https://www.iut-trier.de/livestream/livestreamaac.mp3
	source := &domain.Source{
		Kind: domain.SourceKindURL,
		ID:   "https://www.iut-trier.de/livestream/livestreamaac.mp3",
	}

	event := &recording.RequestedEvent{
		RecordingID: "test-id",
		Source:      source,
		Output:      outputFile,
		Duration:    2 * time.Second, // 短い録音時間
	}

	err := recorder.Rec(context.Background(), event)
	assert.NoError(t, err)
	
	// 出力ファイルが存在するか確認
	_, err = os.Stat(outputFile)
	assert.NoError(t, err)
	require.FileExists(t, outputFile)
}