package recording

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/sun-yryr/recoto/internal/broker"
	"github.com/sun-yryr/recoto/internal/event/common"
	"go.uber.org/zap"
)

// StartedSubject is the NATS subject for recording started events.
var StartedSubject = "recoto.recording.started.v1"

// StartedEvent represents an event that is published when recording has started.
type StartedEvent struct {
	common.CommonEvent
	RecordingID string `json:"recordingId"`
}

// StartedService handles publishing and subscribing to recording started events.
type StartedService struct {
	broker broker.Broker
	logger *zap.Logger
}

// NewStartedService creates a new instance of StartedService.
func NewStartedService(broker broker.Broker, logger *zap.Logger) *StartedService {
	return &StartedService{
		broker: broker,
		logger: logger,
	}
}

// Publish publishes a recording started event to the broker.
func (s *StartedService) Publish(ctx context.Context, recordingID string) error {
	event := StartedEvent{
		CommonEvent: common.CommonEvent{
			Timestamp: time.Now(),
		},
		RecordingID: recordingID,
	}
	message, err := json.Marshal(event)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal %s event", StartedSubject)
	}
	return s.broker.Publish(ctx, StartedSubject, message)
}

// Subscribe registers a handler for recording started events.
func (s *StartedService) Subscribe(ctx context.Context, handler func(context.Context, *StartedEvent)) (broker.UnsubscribeFunc, error) {
	return s.broker.Subscribe(ctx, StartedSubject, func(msg []byte) {
		var event StartedEvent
		if err := json.Unmarshal(msg, &event); err != nil {
			s.logger.Error("failed to unmarshal message", zap.Error(err), zap.String("subject", StartedSubject))
			return
		}
		handler(ctx, &event)
	})
}
