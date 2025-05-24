// Package event はイベントに共通するメタデータ構造を提供する
package event

import (
	"time"

	"github.com/google/uuid"

	"github.com/sun-yryr/recoto/internal/trace"
)

// Metadata はすべてのイベントに共通するメタデータ。
type Metadata struct {
	TraceID   trace.ID  `json:"traceId"`
	Timestamp time.Time `json:"timestamp"`
	EventID   string    `json:"eventId"`
}

// NewMetadata は新しいメタデータを作成する。
func NewMetadata(traceID trace.ID) Metadata {
	return Metadata{
		TraceID:   traceID,
		Timestamp: time.Now(),
		EventID:   uuid.New().String(),
	}
}
