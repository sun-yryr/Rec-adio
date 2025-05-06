package recorder

import (
	"context"
	"time"

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
	r.logger.Info("start recording", zap.String("url", event.Source.ID))
	// TODO: 録音を行う
	time.Sleep(10 * time.Second)
	r.logger.Info("finish recording", zap.String("url", event.Source.ID))

	return nil
}
