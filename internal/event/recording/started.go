package recording

import (
	"time"

	"github.com/sun-yryr/recoto/internal/broker"
	"github.com/sun-yryr/recoto/internal/eventutil"
)

// StartedSubject は録画開始イベントの NATS サブジェクト。
const StartedSubject = "recoto.recording.started.v1"

// StartedEvent は録画開始イベントを表すイベント。
type StartedEvent struct {
	RecordingID string    `json:"recordingId"`
	Timestamp   time.Time `json:"timestamp"`
}

// NewStartedService は StartedEvent 用のサービスを生成するヘルパー。
func NewStartedService(
	broker broker.Broker,
) *eventutil.EventService[StartedEvent] {
	return eventutil.NewEventService[StartedEvent](broker, StartedSubject)
}
