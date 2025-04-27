package nats

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/sun-yryr/recoto/internal/broker"
)

type natsBroker struct {
	nc *nats.Conn
	ns *server.Server
}

func NewEmbeddedBroker() (broker.Broker, error) {
	ns, err := startEmbeddedServer()
	if err != nil {
		return nil, errors.Wrap(err, "failed to start embedded nats server")
	}
	nc, err := nats.Connect(ns.ClientURL())
	if err != nil {
		ns.Shutdown()
		return nil, errors.Wrap(err, "failed to connect to nats server")
	}

	return &natsBroker{
		nc: nc,
		ns: ns,
	}, nil
}

func (b *natsBroker) Publish(ctx context.Context, subject string, message []byte) error {
	return b.nc.Publish(subject, message)
}

func (b *natsBroker) Subscribe(ctx context.Context, subject string, handler func(message []byte)) (broker.UnsubscribeFunc, error) {
	subscription, err := b.nc.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg.Data)
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to subscribe to subject")
	}
	return subscription.Unsubscribe, nil
}

func (b *natsBroker) Close() error {
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
