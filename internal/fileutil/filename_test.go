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
			want:     "_filename_with_spaces_",
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
			want:     "file_name_with_spaces",
		},
		{
			name:     "filename with multiple consecutive underscores",
			filename: "file__name___with____underscores",
			want:     "file_name_with_underscores",
		},
		{
			name:     "Windows reserved name - CON",
			filename: "CON",
			want:     "recording",
		},
		{
			name:     "Windows reserved name - con (lowercase)",
			filename: "con",
			want:     "recording",
		},
		{
			name:     "Windows reserved name - PRN",
			filename: "PRN",
			want:     "recording",
		},
		{
			name:     "Windows reserved name - COM1",
			filename: "COM1",
			want:     "recording",
		},
		{
			name:     "Windows reserved name - LPT9",
			filename: "LPT9",
			want:     "recording",
		},
		{
			name:     "special directory name - dot",
			filename: ".",
			want:     "recording",
		},
		{
			name:     "special directory name - double dot",
			filename: "..",
			want:     "recording",
		},
		{
			name:     "valid filename with reserved substring",
			filename: "CONtent",
			want:     "CONtent",
		},
		{
			name:     "filename with spaces becoming reserved after sanitization",
			filename: "C O N",
			want:     "C_O_N",
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
