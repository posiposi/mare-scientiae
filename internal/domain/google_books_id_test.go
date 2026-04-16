package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNewGoogleBooksID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{
			name:    "有効なID",
			input:   "googleId123",
			want:    "googleId123",
			wantErr: nil,
		},
		{
			name:    "50文字ちょうど",
			input:   strings.Repeat("a", 50),
			want:    strings.Repeat("a", 50),
			wantErr: nil,
		},
		{
			name:    "空文字",
			input:   "",
			want:    "",
			wantErr: ErrGoogleBooksIDRequired,
		},
		{
			name:    "51文字で上限超過",
			input:   strings.Repeat("a", 51),
			want:    "",
			wantErr: ErrGoogleBooksIDTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGoogleBooksID(tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("NewGoogleBooksID(%q) error = %v, want %v", tt.input, err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("NewGoogleBooksID(%q) unexpected error: %v", tt.input, err)
			}
			if got.String() != tt.want {
				t.Errorf("NewGoogleBooksID(%q).String() = %q, want %q", tt.input, got.String(), tt.want)
			}
		})
	}
}
