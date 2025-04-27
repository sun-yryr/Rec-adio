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

var FinishedSubject = "recoto.recording.finished.v1"

type FinishedEvent struct {
	common.CommonEvent
	RecordingID string `json:"recording_id"`
}

type FinishedService struct {
	broker broker.Broker
	logger *zap.Logger
}

func NewFinishedService(broker broker.Broker, logger *zap.Logger) *FinishedService {
	return &FinishedService{
		broker: broker,
		logger: logger,
	}
}

func (s *FinishedService) Publish(ctx context.Context, recordingID string) error {
	e := FinishedEvent{
		CommonEvent: common.CommonEvent{
			Timestamp: time.Now(),
		},
		RecordingID: recordingID,
	}
	message, err := json.Marshal(e)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal %s event", FinishedSubject)
	}
	return s.broker.Publish(ctx, FinishedSubject, message)
}

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
