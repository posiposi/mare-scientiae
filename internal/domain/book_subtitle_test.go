package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNewBookSubtitle(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{
			name:    "有効な副題",
			input:   "副題テスト",
			want:    "副題テスト",
			wantErr: nil,
		},
		{
			name:    "500文字ちょうど",
			input:   strings.Repeat("あ", 500),
			want:    strings.Repeat("あ", 500),
			wantErr: nil,
		},
		{
			name:    "空文字",
			input:   "",
			want:    "",
			wantErr: ErrBookSubtitleEmpty,
		},
		{
			name:    "501文字で上限超過",
			input:   strings.Repeat("あ", 501),
			want:    "",
			wantErr: ErrBookSubtitleTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBookSubtitle(tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("NewBookSubtitle(%q) error = %v, want %v", tt.input, err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("NewBookSubtitle(%q) unexpected error: %v", tt.input, err)
			}
			if got.String() != tt.want {
				t.Errorf("NewBookSubtitle(%q).String() = %q, want %q", tt.input, got.String(), tt.want)
			}
		})
	}
}
