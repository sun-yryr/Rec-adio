// Package broker は、メッセージブローカーのインターフェースを定義します。
package broker

import "context"

// UnsubscribeFunc は、サブスクリプションを解除する関数です。
type UnsubscribeFunc func() error

// Broker は、アプリケーションのメッセージブローカーのインターフェースを定義します。
type Broker interface {
	Publish(ctx context.Context, subject string, message []byte) error
	Subscribe(
		ctx context.Context,
		subject string,
		handler func(message []byte),
	) (UnsubscribeFunc, error)
	Close() error
}
