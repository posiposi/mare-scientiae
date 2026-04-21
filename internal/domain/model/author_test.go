package model

import (
	"errors"
	"testing"
)

func TestNewAuthor(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{
			name:    "有効な著者名",
			input:   "著者A",
			want:    "著者A",
			wantErr: nil,
		},
		{
			name:    "空文字",
			input:   "",
			want:    "",
			wantErr: ErrAuthorEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAuthor(tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("NewAuthor(%q) error = %v, want %v", tt.input, err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("NewAuthor(%q) unexpected error: %v", tt.input, err)
			}
			if got.String() != tt.want {
				t.Errorf("NewAuthor(%q).String() = %q, want %q", tt.input, got.String(), tt.want)
			}
		})
	}
}
