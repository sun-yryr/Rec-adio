// Package recording は録画イベントのパッケージ
package recording

import (
	"time"

	"go.uber.org/zap"

	"github.com/sun-yryr/recoto/internal/broker"
	"github.com/sun-yryr/recoto/internal/eventutil"
)

// FinishedSubject は録画終了イベントの NATS サブジェクト。
const FinishedSubject = "recoto.recording.finished.v1"

// FinishedEvent は録画終了イベントを表すイベント。
type FinishedEvent struct {
	RecordingID string    `json:"recordingId"`
	Success     bool      `json:"success"`
	Error       *string   `json:"error,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}

// NewFinishedService は FinishedEvent 用のサービスを生成するヘルパー。
func NewFinishedService(
	broker broker.Broker,
	logger *zap.Logger,
) *eventutil.EventService[FinishedEvent] {
	return eventutil.NewEventService[FinishedEvent](broker, logger, FinishedSubject)
}
