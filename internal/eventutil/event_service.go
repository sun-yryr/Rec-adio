// Package eventutil はイベントの発行・購読を行う関数を提供する
package eventutil

import (
	"context"
	"encoding/json"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"

	"github.com/sun-yryr/recoto/internal/broker"
)

// EventService はイベントを発行・購読する関数を提供するジェネリックなサービス。
type EventService[T any] struct {
	broker  broker.Broker
	logger  *zap.Logger
	subject string
}

// NewEventService は新しい EventService を作成する。
func NewEventService[T any](
	broker broker.Broker,
	logger *zap.Logger,
	subject string,
) *EventService[T] {
	return &EventService[T]{broker: broker, logger: logger, subject: subject}
}

// Publish はイベントを発行する。
func (s *EventService[T]) Publish(ctx context.Context, event T) error {
	message, err := json.Marshal(event)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal %s event", s.subject)
	}

	if err := s.broker.Publish(ctx, s.subject, message); err != nil {
		return errors.Wrapf(err, "failed to publish %s event", s.subject)
	}

	return nil
}

// Subscribe はイベントを購読する。
func (s *EventService[T]) Subscribe(
	ctx context.Context,
	handler func(context.Context, *T),
) (broker.UnsubscribeFunc, error) {
	subscribeHandler := func(msg []byte) {
		var event T
		if err := json.Unmarshal(msg, &event); err != nil {
			s.logger.Error(
				"failed to unmarshal message",
				zap.Error(err),
				zap.String("subject", s.subject),
			)

			return
		}

		handler(ctx, &event)
	}

	unsubscribe, err := s.broker.Subscribe(ctx, s.subject, subscribeHandler)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to subscribe to %s", s.subject)
	}

	return unsubscribe, nil
}
