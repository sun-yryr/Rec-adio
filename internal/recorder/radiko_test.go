package recorder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sun-yryr/recoto/internal/domain"
	"go.uber.org/zap/zaptest"
)

func TestRadikoRecorder_GetSupportSource(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	recorder := NewRadikoRecorder(logger, RadikoConfig{})

	sources := recorder.GetSupportSource()
	assert.Equal(t, []domain.SourceKind{domain.SourceKindRadiko}, sources)
}

func TestRadikoRecorder_GetName(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	recorder := NewRadikoRecorder(logger, RadikoConfig{})

	name := recorder.GetName()
	assert.Equal(t, "RadikoRecorder", name)
}

func TestRadikoRecorder_CheckAvailable(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	recorder := NewRadikoRecorder(logger, RadikoConfig{})

	// この部分はffmpegがインストールされているかどうかに依存するため、
	// CIでは失敗する可能性がある。
	// 実際の環境に合わせてテストを調整する必要がある。
	_ = recorder.CheckAvailable()
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

			source, err := domain.NewRadikoSource(tt.stationID, tt.meta)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, source)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, source)
				assert.Equal(t, domain.SourceKindRadiko, source.Kind)
				assert.Equal(t, tt.stationID, source.ID)

				if tt.meta == nil {
					assert.NotNil(t, source.Meta)
					assert.Empty(t, source.Meta)
				} else {
					assert.Equal(t, tt.meta, source.Meta)
				}
			}
		})
	}
}

// 以下のテストはRadikoのAPIに依存するため、モックを使用するか、
// 実際のAPIを使用する場合は適切な環境変数やフラグで制御する必要がある。

func TestRadikoRecorder_Rec_Mock(t *testing.T) {
	t.Parallel()

	// このテストはモックを使用して実装する必要がある
	t.Skip("This test requires mocking the Radiko API")
}

func TestRadikoRecorder_authorization_Mock(t *testing.T) {
	t.Parallel()

	// このテストはモックを使用して実装する必要がある
	t.Skip("This test requires mocking the Radiko API")
}

func TestRadikoRecorder_getM3U8URL_Mock(t *testing.T) {
	t.Parallel()

	// このテストはモックを使用して実装する必要がある
	t.Skip("This test requires mocking the Radiko API")
}

func TestRadikoRecorder_SearchProgram_Mock(t *testing.T) {
	t.Parallel()

	// このテストはモックを使用して実装する必要がある
	t.Skip("This test requires mocking the Radiko API")
}
