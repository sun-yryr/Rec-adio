package broker

import "context"

type UnsubscribeFunc func() error

type Broker interface {
	Publish(ctx context.Context, subject string, message []byte) error
	Subscribe(ctx context.Context, subject string, handler func(message []byte)) (UnsubscribeFunc, error)
	Close() error
}
