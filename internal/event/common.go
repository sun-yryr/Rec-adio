package event

import "time"

type CommonEvent struct {
	Timestamp time.Time `json:"timestamp"`
}
