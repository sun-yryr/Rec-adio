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

var StartedSubject = "recoto.recording.started.v1"

type StartedEvent struct {
	common.CommonEvent
	RecordingID string `json:"recording_id"`
}

type StartedService struct {
	broker broker.Broker
	logger *zap.Logger
}

func NewStartedService(broker broker.Broker, logger *zap.Logger) *StartedService {
	return &StartedService{
		broker: broker,
		logger: logger,
	}
}

func (s *StartedService) Publish(ctx context.Context, recordingID string) error {
	e := StartedEvent{
		CommonEvent: common.CommonEvent{
			Timestamp: time.Now(),
		},
		RecordingID: recordingID,
	}
	message, err := json.Marshal(e)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal %s event", StartedSubject)
	}
	return s.broker.Publish(ctx, StartedSubject, message)
}

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
