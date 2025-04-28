package broker

import "context"

// UnsubscribeFunc is a function to unsubscribe from a subscription.
type UnsubscribeFunc func() error

// Broker defines the interface for message brokers in the application.
type Broker interface {
	Publish(ctx context.Context, subject string, message []byte) error
	Subscribe(ctx context.Context, subject string, handler func(message []byte)) (UnsubscribeFunc, error)
	Close() error
}
