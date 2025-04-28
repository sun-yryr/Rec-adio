package common

import "time"

// CommonEvent contains common fields for all events.
type CommonEvent struct {
	Timestamp time.Time `json:"timestamp"`
}
