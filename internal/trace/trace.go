// Package trace は分散トレーシング用のIDを管理する
package trace

import (
	"context"

	"github.com/google/uuid"
)

type ctxKeyTraceID struct{}

// ID は分散トレーシング用の一意識別子。
type ID string

// NewTraceID は新しいIDを生成する。
func NewTraceID() ID {
	return ID(uuid.Must(uuid.NewV7()).String())
}

// WithTraceID は context にIDを追加する。
func WithTraceID(ctx context.Context, traceID ID) context.Context {
	return context.WithValue(ctx, ctxKeyTraceID{}, traceID)
}

// IDFromContext は context からIDを取得する。
func IDFromContext(ctx context.Context) (ID, bool) {
	traceID, ok := ctx.Value(ctxKeyTraceID{}).(ID)

	return traceID, ok
}

// GetOrGenerate は context からIDを取得し、存在しない場合は生成する。
func GetOrGenerate(ctx context.Context) (context.Context, ID) {
	if traceID, ok := IDFromContext(ctx); ok {
		return ctx, traceID
	}

	traceID := NewTraceID()

	return WithTraceID(ctx, traceID), traceID
}

// String はIDを文字列として返す。
func (t ID) String() string {
	return string(t)
}
