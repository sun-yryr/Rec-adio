// Package broker は、メッセージブローカーのインターフェースを定義する
package broker

import "context"

// UnsubscribeFunc は、サブスクリプションを解除する関数。
type UnsubscribeFunc func() error

// Broker は、アプリケーションのメッセージブローカーのインターフェースを定義する。
type Broker interface {
	Publish(ctx context.Context, subject string, message interface{}) error
	Subscribe(
		ctx context.Context,
		subject string,
		handler func(context.Context, []byte),
	) (UnsubscribeFunc, error)
	Close() error
}
