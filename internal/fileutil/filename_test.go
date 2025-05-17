package fileutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeFilename(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "normal filename",
			filename: "normal_filename",
			want:     "normal_filename",
		},
		{
			name:     "filename with invalid chars",
			filename: "file/with\\invalid:chars*?\"<>|",
			want:     "file_with_invalid_chars_",
		},
		{
			name:     "filename with consecutive invalid chars",
			filename: "file//with\\\\invalid::chars",
			want:     "file_with_invalid_chars",
		},
		{
			name:     "filename with spaces",
			filename: "  filename with spaces  ",
			want:     "filename with spaces",
		},
		{
			name:     "empty filename",
			filename: "",
			want:     "recording",
		},
		{
			name:     "filename with only invalid chars",
			filename: "/:*?\"<>|",
			want:     "_",
		},
		{
			name:     "filename with spaces and invalid chars",
			filename: "  //:*?\"<>|  ",
			want:     "_",
		},
		{
			name:     "filename with multiple consecutive spaces",
			filename: "file  name  with    spaces",
			want:     "file  name  with    spaces",
		},
		{
			name:     "filename with multiple consecutive underscores",
			filename: "file__name___with____underscores",
			want:     "file_name_with_underscores",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := SanitizeFilename(tt.filename)
			assert.Equal(t, tt.want, got)
		})
	}
}
