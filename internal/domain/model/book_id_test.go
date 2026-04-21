package model

import (
	"errors"
	"testing"
)

func TestNewBookID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{
			name:    "有効なUUID",
			input:   "550e8400-e29b-41d4-a716-446655440000",
			want:    "550e8400-e29b-41d4-a716-446655440000",
			wantErr: nil,
		},
		{
			name:    "空文字",
			input:   "",
			want:    "",
			wantErr: ErrBookIDRequired,
		},
		{
			name:    "不正なUUID形式",
			input:   "not-a-uuid",
			want:    "",
			wantErr: ErrBookIDInvalidFormat,
		},
		{
			name:    "大文字を含むUUID",
			input:   "550E8400-E29B-41D4-A716-446655440000",
			want:    "",
			wantErr: ErrBookIDInvalidFormat,
		},
		{
			name:    "ハイフンなしのUUID",
			input:   "550e8400e29b41d4a716446655440000",
			want:    "",
			wantErr: ErrBookIDInvalidFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBookID(tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("NewBookID(%q) error = %v, want %v", tt.input, err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("NewBookID(%q) unexpected error: %v", tt.input, err)
			}
			if got.String() != tt.want {
				t.Errorf("NewBookID(%q).String() = %q, want %q", tt.input, got.String(), tt.want)
			}
		})
	}
}
