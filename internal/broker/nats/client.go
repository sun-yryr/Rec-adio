// Package nats は、NATSを使用したブローカーの実装を提供する
package nats

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/sun-yryr/recoto/internal/broker"
	"github.com/sun-yryr/recoto/internal/event"
	"github.com/sun-yryr/recoto/internal/logger"
	"github.com/sun-yryr/recoto/internal/trace"
)

// Message は TraceID を含むメッセージ構造。
type Message struct {
	Metadata event.Metadata  `json:"metadata"`
	Payload  json.RawMessage `json:"payload"`
}

// Broker は、NATSのブローカーを表す構造体。
type Broker struct {
	nc *nats.Conn
	ns *server.Server
}

// NewEmbeddedBroker は、NATSのブローカーを作成する。
func NewEmbeddedBroker(logger *zap.Logger) (*Broker, error) {
	server, err := startEmbeddedServer()
	if err != nil {
		return nil, errors.Wrap(err, "failed to start embedded nats server")
	}

	conn, err := nats.Connect(server.ClientURL())
	if err != nil {
		server.Shutdown()

		return nil, errors.Wrap(err, "failed to connect to nats server")
	}

	logger.Debug("connected to nats server", zap.String("client_url", server.ClientURL()))

	return &Broker{
		nc: conn,
		ns: server,
	}, nil
}

// Publish は context を考慮してメッセージを送信する。
func (b *Broker) Publish(ctx context.Context, subject string, payload interface{}) error {
	// context からTraceIDを取得または生成
	ctx, traceID := trace.GetOrGenerate(ctx)
	log := logger.FromContextWithTrace(ctx)

	// ペイロードをJSON化
	payloadData, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "failed to marshal payload")
	}

	// 拡張メッセージの作成
	msg := Message{
		Metadata: event.NewMetadata(traceID),
		Payload:  payloadData,
	}

	msgData, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "failed to marshal message")
	}

	// context のデッドラインを考慮
	if deadline, ok := ctx.Deadline(); ok {
		timeout := time.Until(deadline)
		if timeout <= 0 {
			return context.DeadlineExceeded
		}

		// タイムアウト付きで送信
		done := make(chan error, 1)
		go func() {
			done <- b.nc.Publish(subject, msgData)
		}()

		select {
		case err := <-done:
			if err != nil {
				return errors.Wrap(err, "failed to publish message")
			}

			log.Debug("published message", zap.String("subject", subject))

			return nil
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), "context canceled while publishing message")
		case <-time.After(timeout):
			return context.DeadlineExceeded
		}
	}

	// 通常の送信
	if err := b.nc.Publish(subject, msgData); err != nil {
		return errors.Wrap(err, "failed to publish message")
	}

	log.Debug("published message", zap.String("subject", subject))

	return nil
}

// Subscribe は context を考慮してサブスクリプションを作成する。
func (b *Broker) Subscribe(
	ctx context.Context,
	subject string,
	handler func(context.Context, []byte),
) (broker.UnsubscribeFunc, error) {
	baseLogger := logger.FromContextWithTrace(ctx)

	subscription, err := b.nc.Subscribe(subject, func(msg *nats.Msg) {
		// 拡張メッセージをデコード
		var enhancedMsg Message
		if err := json.Unmarshal(msg.Data, &enhancedMsg); err != nil {
			baseLogger.Error("failed to unmarshal enhanced message", zap.Error(err))

			return
		}

		// TraceID付きの新しい context を作成
		msgCtx := trace.WithTraceID(ctx, enhancedMsg.Metadata.TraceID)

		// ロガーを拡張
		msgLogger := logger.FromContextWithTrace(msgCtx)
		msgCtx = logger.WithLogger(msgCtx, msgLogger)

		msgLogger.Debug("received message", zap.String("subject", msg.Subject))

		// ハンドラーを実行
		handler(msgCtx, enhancedMsg.Payload)
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to subscribe to subject")
	}

	baseLogger.Debug("subscribed to subject", zap.String("subject", subject))

	// context のキャンセレーションでサブスクリプションを終了
	go func() {
		<-ctx.Done()

		if err := subscription.Unsubscribe(); err != nil {
			baseLogger.Error(
				"failed to unsubscribe from subject",
				zap.Error(err),
				zap.String("subject", subject),
			)
		}
	}()

	return subscription.Unsubscribe, nil
}

// Close は、NATSのブローカーを閉じる。
func (b *Broker) Close() error {
	var closeErr error

	if err := b.nc.Drain(); err != nil {
		closeErr = errors.Wrap(err, "failed to drain nats connection")
	}

	b.nc.Close()

	if b.ns != nil {
		b.ns.Shutdown()
	}

	return closeErr
}
