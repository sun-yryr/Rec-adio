package recorder

import (
	"context"
	"strconv"
	"time"

	"github.com/cockroachdb/errors"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"go.uber.org/zap"

	"github.com/sun-yryr/recoto/internal/domain"
	"github.com/sun-yryr/recoto/internal/event/recording"
)

// URLRecorder はURLをソースとして録音を行うrecorder。
type URLRecorder struct {
	logger *zap.Logger
}

// NewURLRecorder はURLRecorderを生成するコンストラクタ。
func NewURLRecorder(logger *zap.Logger) *URLRecorder {
	return &URLRecorder{
		logger: logger,
	}
}

// GetSupportSource はURLRecorderがサポートするソースを返す。
func (r *URLRecorder) GetSupportSource() []domain.SourceKind {
	return []domain.SourceKind{domain.SourceKindURL}
}

// Rec はURLをソースとして録音を行う。
func (r *URLRecorder) Rec(ctx context.Context, event *recording.RequestedEvent) error {
	cmd := ffmpeg_go.
		Input(event.Source.ID, ffmpeg_go.KwArgs{"t": strconv.FormatInt(event.Duration, 10)}).
		Audio().
		Output(event.Output, ffmpeg_go.KwArgs{"acodec": "aac"})

	// キャンセル付きのcontextを設定する
	cmd.Context = ctx
	// 処理の猶予時間を追加
	const gracePeriod = 10 * time.Second
	cmd = cmd.WithTimeout(time.Duration(event.Duration)*time.Second + gracePeriod)

	r.logger.Debug(
		"start recording",
		zap.String("recordingId", event.RecordingID),
		zap.String("cmd", cmd.String()),
	)

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to run ffmpeg")
	}

	r.logger.Debug(
		"finish recording",
		zap.String("recordingId", event.RecordingID),
	)

	return nil
}
