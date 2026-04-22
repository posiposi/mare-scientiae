package model

import (
	"errors"
	"strings"
	"testing"
)

func TestNewBookTitle(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{
			name:    "有効なタイトル",
			input:   "テストタイトル",
			want:    "テストタイトル",
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
			wantErr: ErrBookTitleRequired,
		},
		{
			name:    "501文字で上限超過",
			input:   strings.Repeat("あ", 501),
			want:    "",
			wantErr: ErrBookTitleTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBookTitle(tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("NewBookTitle(%q) error = %v, want %v", tt.input, err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("NewBookTitle(%q) unexpected error: %v", tt.input, err)
			}
			if got.String() != tt.want {
				t.Errorf("NewBookTitle(%q).String() = %q, want %q", tt.input, got.String(), tt.want)
			}
		})
	}
}
