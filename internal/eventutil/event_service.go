// Package eventutil はイベントの発行・購読を行う関数を提供する
package eventutil

import (
	"context"
	"encoding/json"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"

	"github.com/sun-yryr/recoto/internal/broker"
	"github.com/sun-yryr/recoto/internal/logger"
)

// EventService はイベントを発行・購読する関数を提供するジェネリックなサービス。
type EventService[T any] struct {
	broker  broker.Broker
	subject string
}

// NewEventService は新しい EventService を作成する。
func NewEventService[T any](
	broker broker.Broker,
	subject string,
) *EventService[T] {
	return &EventService[T]{broker: broker, subject: subject}
}

// EncodeEvent はイベントをJSON形式にエンコードする（テスト用）。
func EncodeEvent[T any](event T) ([]byte, error) {
	data, err := json.Marshal(event)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal event")
	}

	return data, nil
}

// Publish は TraceID 付きでイベントを発行する。
func (s *EventService[T]) Publish(ctx context.Context, event T) error {
	// context からロガーを取得して TraceID で拡張
	log := logger.FromContextWithTrace(ctx)
	log.Debug("publishing event", zap.String("subject", s.subject))

	if err := s.broker.Publish(ctx, s.subject, event); err != nil {
		return errors.Wrapf(err, "failed to publish %s event", s.subject)
	}

	return nil
}

// Subscribe は TraceID 対応でイベントを購読する。
func (s *EventService[T]) Subscribe(
	ctx context.Context,
	handler func(context.Context, *T),
) (broker.UnsubscribeFunc, error) {
	subscribeHandler := func(msgCtx context.Context, msg []byte) {
		var event T
		if err := json.Unmarshal(msg, &event); err != nil {
			log := logger.FromContextWithTrace(msgCtx)
			log.Error(
				"failed to unmarshal message",
				zap.Error(err),
				zap.String("subject", s.subject),
			)

			return
		}

		// TraceID 付きの context でハンドラーを実行
		handler(msgCtx, &event)
	}

	unsubscribe, err := s.broker.Subscribe(ctx, s.subject, subscribeHandler)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to subscribe to %s", s.subject)
	}

	return unsubscribe, nil
}
