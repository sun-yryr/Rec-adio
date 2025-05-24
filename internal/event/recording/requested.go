package recording

import (
	"time"

	"github.com/sun-yryr/recoto/internal/broker"
	"github.com/sun-yryr/recoto/internal/domain"
	"github.com/sun-yryr/recoto/internal/eventutil"
)

// RequestedSubject は録画要求イベントの NATS サブジェクト。
const RequestedSubject = "recoto.recording.requested.v1"

// RequestedEvent は録画要求イベントを表すイベント。
type RequestedEvent struct {
	RecordingID string                 `json:"recordingId"`
	Status      domain.RecordingStatus `json:"status"`
	Source      *domain.Source         `json:"source"`
	Output      string                 `json:"output"`
	Duration    time.Duration          `json:"duration"`
	Timestamp   time.Time              `json:"timestamp"`
}

// NewRequestedEvent はRecordingからRequestedEventを生成する。
func NewRequestedEvent(from domain.Recording) *RequestedEvent {
	return &RequestedEvent{
		RecordingID: string(from.ID),
		Status:      from.Status,
		Source:      from.Source,
		Output:      from.Output,
		Duration:    from.Duration,
		Timestamp:   time.Now(),
	}
}

// ToDomain はRequestedEventからRecordingを生成する。
func (e *RequestedEvent) ToDomain() *domain.Recording {
	return &domain.Recording{
		ID:       domain.RecordingID(e.RecordingID),
		Status:   e.Status,
		Source:   e.Source,
		Output:   e.Output,
		Duration: e.Duration,
	}
}

// NewRequestedService は RequestedEvent 用のサービスを生成するヘルパー。
func NewRequestedService(
	broker broker.Broker,
) *eventutil.EventService[RequestedEvent] {
	return eventutil.NewEventService[RequestedEvent](broker, RequestedSubject)
}
