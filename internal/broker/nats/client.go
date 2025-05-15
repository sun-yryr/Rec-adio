// Package nats は、NATSを使用したブローカーの実装を提供する
package nats

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/sun-yryr/recoto/internal/broker"
)

// Broker は、NATSのブローカーを表す構造体。
type Broker struct {
	nc     *nats.Conn
	ns     *server.Server
	logger *zap.Logger
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
		nc:     conn,
		ns:     server,
		logger: logger,
	}, nil
}

// Publish は、NATSのブローカーにメッセージを送信する。
func (b *Broker) Publish(_ context.Context, subject string, message []byte) error {
	if err := b.nc.Publish(subject, message); err != nil {
		return errors.Wrap(err, "failed to publish message")
	}

	b.logger.Debug("published message", zap.String("subject", subject))

	return nil
}

// Subscribe は、NATSのブローカーにメッセージを受信するサブスクリプションを作成する。
func (b *Broker) Subscribe(
	_ context.Context,
	subject string,
	handler func(message []byte),
) (broker.UnsubscribeFunc, error) {
	subscription, err := b.nc.Subscribe(subject, func(msg *nats.Msg) {
		b.logger.Debug(
			"received message",
			zap.String("subject", msg.Subject),
			zap.String("data", string(msg.Data)),
		)
		handler(msg.Data)
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to subscribe to subject")
	}

	b.logger.Debug("subscribed to subject", zap.String("subject", subject))

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
