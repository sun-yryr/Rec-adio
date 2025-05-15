// Package domain は、アプリケーションのドメイン層を定義します.
package domain

import (
	stdurl "net/url"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

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
func NewURLSource(url string) (*Source, error) {
	if _, err := stdurl.Parse(url); err != nil {
		return nil, errors.Wrap(err, "invalid url")
	}

	return &Source{
		Kind: SourceKindURL,
		ID:   url,
		Meta: make(map[string]string),
	}, nil
}

// Recording は、録音を表す型.
type Recording struct {
	ID       RecordingID     `json:"id"`
	Source   *Source         `json:"source"`
	Status   RecordingStatus `json:"status"`
	Output   string          `json:"output"`
	Duration time.Duration   `json:"duration"`
}

// NewRecording は、新しい録音を作成する.
func NewRecording(source *Source, output string, duration time.Duration) (*Recording, error) {
	if source == nil {
		return nil, errors.New("source is nil")
	}

	if output == "" {
		return nil, errors.New("output is empty")
	}

	if duration <= 0 {
		return nil, errors.New("duration must be positive")
	}

	return &Recording{
		ID:       RecordingID(uuid.New().String()),
		Source:   source,
		Status:   RecordingStatusRequested,
		Output:   output,
		Duration: duration,
	}, nil
}

// ErrInvalidRecordingStatus は録音ステータスが無効な場合のエラー。
var ErrInvalidRecordingStatus = errors.New("invalid recording status")

// Start は録音の状態を実行中に変更する。
func (r *Recording) Start() error {
	if r.Status != RecordingStatusRequested {
		return errors.Wrapf(
			ErrInvalidRecordingStatus,
			"cannot start recording with status: %s",
			r.Status,
		)
	}

	r.Status = RecordingStatusRunning

	return nil
}

// Finish は録音の状態を完了に変更する。
func (r *Recording) Finish() error {
	if r.Status != RecordingStatusRunning {
		return errors.Wrapf(
			ErrInvalidRecordingStatus,
			"cannot finish recording with status: %s",
			r.Status,
		)
	}

	r.Status = RecordingStatusFinished

	return nil
}

// Fail は録音の状態を失敗に変更する。
func (r *Recording) Fail() error {
	r.Status = RecordingStatusFailed

	return nil
}
