// Package domain は、アプリケーションのドメイン層を定義します.
package domain

import "github.com/google/uuid"

// RecordingID は、録音のIDを表す型.
type RecordingID string

// RecordingStatus は、録音のステータスを表す型.
type RecordingStatus string

const (
	// RecordingStatusRequested は、録音が要求された状態を表す.
	RecordingStatusRequested RecordingStatus = "requested"
	// RecordingStatusRunning は、録音が実行中の状態を表す.
	RecordingStatusRunning RecordingStatus = "running"
	// RecordingStatusFinished は、録音が完了した状態を表す.
	RecordingStatusFinished RecordingStatus = "finished"
	// RecordingStatusFailed は、録音が失敗した状態を表す.
	RecordingStatusFailed RecordingStatus = "failed"
)

// SourceKind は、録音のソースの種類を表す型.
type SourceKind string

const (
	// SourceKindURL は、URLのソースを表す.
	SourceKindURL SourceKind = "url"
)

// Source は、録音のソースを表す型.
type Source struct {
	Kind SourceKind        `json:"kind"`
	ID   string            `json:"id"`
	Meta map[string]string `json:"meta"`
}

// NewURLSource は、URLのソースを作成する.
func NewURLSource(url string) Source {
	return Source{
		Kind: SourceKindURL,
		ID:   url,
	}
}

// Recording は、録音を表す型.
type Recording struct {
	ID       RecordingID     `json:"id"`
	Source   Source          `json:"source"`
	Status   RecordingStatus `json:"status"`
	Output   string          `json:"output"`
	Duration int64           `json:"duration"`
}

// NewRecording は、新しい録音を作成する.
func NewRecording(source Source, output string, duration int64) *Recording {
	return &Recording{
		ID:       RecordingID(uuid.New().String()),
		Source:   source,
		Status:   RecordingStatusRequested,
		Output:   output,
		Duration: duration,
	}
}
