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

// FinishedSubject is the NATS subject for recording finished events.
var FinishedSubject = "recoto.recording.finished.v1"

// FinishedEvent represents an event that is published when recording has finished.
type FinishedEvent struct {
	common.CommonEvent
	RecordingID string `json:"recording_id"`
}

// FinishedService handles publishing and subscribing to recording finished events.
type FinishedService struct {
	broker broker.Broker
	logger *zap.Logger
}

// NewFinishedService creates a new instance of FinishedService.
func NewFinishedService(broker broker.Broker, logger *zap.Logger) *FinishedService {
	return &FinishedService{
		broker: broker,
		logger: logger,
	}
}

// Publish publishes a recording finished event to the broker.
func (s *FinishedService) Publish(ctx context.Context, recordingID string) error {
	event := FinishedEvent{
		CommonEvent: common.CommonEvent{
			Timestamp: time.Now(),
		},
		RecordingID: recordingID,
	}
	message, err := json.Marshal(event)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal %s event", FinishedSubject)
	}
	return s.broker.Publish(ctx, FinishedSubject, message)
}

// Subscribe registers a handler for recording finished events.
func (s *FinishedService) Subscribe(ctx context.Context, handler func(context.Context, *FinishedEvent)) (broker.UnsubscribeFunc, error) {
	return s.broker.Subscribe(ctx, FinishedSubject, func(msg []byte) {
		var event FinishedEvent
		if err := json.Unmarshal(msg, &event); err != nil {
			s.logger.Error("failed to unmarshal message", zap.Error(err), zap.String("subject", FinishedSubject))
			return
		}
		handler(ctx, &event)
	})
}
