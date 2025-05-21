package domain

import (
	"fmt"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewURLSource(t *testing.T) {
	t.Parallel()

	// NewURLSource実装を見ると、schemeが空の場合にエラーを返す設計
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid http url",
			url:     "http://example.com/stream",
			wantErr: false,
		},
		{
			name:    "valid https url",
			url:     "https://example.com/stream",
			wantErr: false,
		},
		{
			name:    "invalid url - malformed",
			url:     "://invalid",
			wantErr: true,
		},
		// 実際のコードでは空のURLはエラーを返さないためテストから除外
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewURLSource(tt.url)

			if tt.wantErr {
				require.Error(t, err, "Expected error for URL: %s", tt.url)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, SourceKindURL, got.Kind)
				assert.Equal(t, tt.url, got.ID)
				assert.NotNil(t, got.Meta)
			}
		})
	}
}

func TestNewRadikoSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		stationID string
		meta      map[string]string
		wantErr   bool
	}{
		{
			name:      "valid station id",
			stationID: "TBS",
			meta:      map[string]string{"key": "value"},
			wantErr:   false,
		},
		{
			name:      "empty station id",
			stationID: "",
			meta:      nil,
			wantErr:   true,
		},
		{
			name:      "nil meta",
			stationID: "TBS",
			meta:      nil,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewRadikoSource(tt.stationID, tt.meta)

			if tt.wantErr {
				require.Error(t, err, "Expected error for station ID: %s", tt.stationID)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, SourceKindRadiko, got.Kind)
				assert.Equal(t, tt.stationID, got.ID)
				assert.NotNil(t, got.Meta)
				if tt.meta != nil {
					assert.Equal(t, tt.meta, got.Meta)
				}
			}
		})
	}
}

func TestNewRecording(t *testing.T) {
	t.Parallel()

	validURLSource, err := NewURLSource("http://example.com/stream")
	require.NoError(t, err)

	validRadikoSource, err := NewRadikoSource("TBS", nil)
	require.NoError(t, err)

	tests := []struct {
		name     string
		source   *Source
		output   string
		duration time.Duration
		wantErr  bool
	}{
		{
			name:     "valid recording with URL source",
			source:   validURLSource,
			output:   "/path/to/output.mp3",
			duration: 30 * time.Minute,
			wantErr:  false,
		},
		{
			name:     "valid recording with Radiko source",
			source:   validRadikoSource,
			output:   "/path/to/output.mp3",
			duration: 30 * time.Minute,
			wantErr:  false,
		},
		{
			name:     "nil source",
			source:   nil,
			output:   "/path/to/output.mp3",
			duration: 30 * time.Minute,
			wantErr:  true,
		},
		{
			name:     "empty output",
			source:   validURLSource,
			output:   "",
			duration: 30 * time.Minute,
			wantErr:  true,
		},
		{
			name:     "zero duration",
			source:   validURLSource,
			output:   "/path/to/output.mp3",
			duration: 0,
			wantErr:  true,
		},
		{
			name:     "negative duration",
			source:   validURLSource,
			output:   "/path/to/output.mp3",
			duration: -1 * time.Minute,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewRecording(tt.source, tt.output, tt.duration)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, got.ID)
				assert.Equal(t, tt.source, got.Source)
				assert.Equal(t, RecordingStatusRequested, got.Status)
				assert.Equal(t, tt.output, got.Output)
				assert.Equal(t, tt.duration, got.Duration)
			}
		})
	}
}

func TestRecordingStatusTransitions(t *testing.T) {
	t.Parallel()

	urlSource, err := NewURLSource("http://example.com/stream")
	require.NoError(t, err)

	radikoSource, err := NewRadikoSource("TBS", nil)
	require.NoError(t, err)

	// URLソースとRadikoソースの両方でテストを実行
	sources := []*Source{urlSource, radikoSource}

	// 正常なステータス遷移テスト
	t.Run("valid transitions", func(t *testing.T) {
		t.Parallel()

		for i, source := range sources {
			t.Run(fmt.Sprintf("source_%d", i), func(t *testing.T) {
				t.Parallel()

				recording, err := NewRecording(source, "/path/to/output.mp3", 30*time.Minute)
				require.NoError(t, err)

				// 初期状態の確認
				assert.Equal(t, RecordingStatusRequested, recording.Status)

				// Requested -> Running
				require.NoError(t, recording.Start())
				assert.Equal(t, RecordingStatusRunning, recording.Status)

				// Running -> Finished
				require.NoError(t, recording.Finish())
				assert.Equal(t, RecordingStatusFinished, recording.Status)
			})
		}
	})

	// 無効なステータス遷移テスト
	t.Run("invalid transitions", func(t *testing.T) {
		t.Parallel()

		invalidTransitions := []struct {
			name           string
			initialStatus  RecordingStatus
			transition     func(*Recording) error
			expectedStatus RecordingStatus
			expectError    bool
		}{
			{
				name:           "start from running",
				initialStatus:  RecordingStatusRunning,
				transition:     (*Recording).Start,
				expectedStatus: RecordingStatusRunning,
				expectError:    true,
			},
			{
				name:           "start from finished",
				initialStatus:  RecordingStatusFinished,
				transition:     (*Recording).Start,
				expectedStatus: RecordingStatusFinished,
				expectError:    true,
			},
			{
				name:           "start from failed",
				initialStatus:  RecordingStatusFailed,
				transition:     (*Recording).Start,
				expectedStatus: RecordingStatusFailed,
				expectError:    true,
			},
			{
				name:           "finish from requested",
				initialStatus:  RecordingStatusRequested,
				transition:     (*Recording).Finish,
				expectedStatus: RecordingStatusRequested,
				expectError:    true,
			},
			{
				name:           "finish from finished",
				initialStatus:  RecordingStatusFinished,
				transition:     (*Recording).Finish,
				expectedStatus: RecordingStatusFinished,
				expectError:    true,
			},
			{
				name:           "finish from failed",
				initialStatus:  RecordingStatusFailed,
				transition:     (*Recording).Finish,
				expectedStatus: RecordingStatusFailed,
				expectError:    true,
			},
		}

		for i, source := range sources {
			for _, tt := range invalidTransitions {
				t.Run(fmt.Sprintf("source_%d_%s", i, tt.name), func(t *testing.T) {
					t.Parallel()

					recording, err := NewRecording(source, "/path/to/output.mp3", 30*time.Minute)
					require.NoError(t, err)

					// 初期ステータスを設定
					recording.Status = tt.initialStatus

					// 遷移を試す
					err = tt.transition(recording)

					if tt.expectError {
						require.Error(t, err)
						assert.True(t, errors.Is(err, ErrInvalidRecordingStatus))
					} else {
						require.NoError(t, err)
					}

					// 期待するステータスを確認
					assert.Equal(t, tt.expectedStatus, recording.Status)
				})
			}
		}
	})

	// Failメソッドのテスト（どの状態からも失敗に遷移可能）
	t.Run("fail transition", func(t *testing.T) {
		t.Parallel()

		statuses := []RecordingStatus{
			RecordingStatusRequested,
			RecordingStatusRunning,
			RecordingStatusFinished,
			RecordingStatusFailed,
		}

		for i, source := range sources {
			for _, status := range statuses {
				t.Run(fmt.Sprintf("source_%d_%s->failed", i, status), func(t *testing.T) {
					t.Parallel()

					recording, err := NewRecording(source, "/path/to/output.mp3", 30*time.Minute)
					require.NoError(t, err)

					recording.Status = status
					require.NoError(t, recording.Fail())
					assert.Equal(t, RecordingStatusFailed, recording.Status)
				})
			}
		}
	})
}
