package domain

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestNewBookSubtitle(t *testing.T) {
	valid := "副題テスト"
	max := strings.Repeat("あ", 500)
	empty := ""
	tooLong := strings.Repeat("あ", 501)

	tests := []struct {
		name    string
		input   *string
		wantNil bool
		wantStr string
		wantErr error
	}{
		{
			name:    "nilはnilを返す",
			input:   nil,
			wantNil: true,
		},
		{
			name:    "有効な副題",
			input:   &valid,
			wantStr: valid,
		},
		{
			name:    "500文字ちょうど",
			input:   &max,
			wantStr: max,
		},
		{
			name:    "空文字",
			input:   &empty,
			wantErr: ErrBookSubtitleEmpty,
		},
		{
			name:    "501文字で上限超過",
			input:   &tooLong,
			wantErr: ErrBookSubtitleTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBookSubtitle(tt.input)
			input := formatNillableString(tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("NewBookSubtitle(%s) error = nil, want %v", input, tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("NewBookSubtitle(%s) error = %v, want %v", input, err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("NewBookSubtitle(%s) unexpected error: %v", input, err)
			}
			if tt.wantNil {
				if got != nil {
					t.Errorf("NewBookSubtitle(%s) = %q, want nil", input, got.String())
				}
				return
			}
			if got == nil {
				t.Fatalf("NewBookSubtitle(%s) = nil, want %q", input, tt.wantStr)
			}
			if got.String() != tt.wantStr {
				t.Errorf("NewBookSubtitle(%s).String() = %q, want %q", input, got.String(), tt.wantStr)
			}
		})
	}
}

func formatNillableString(v *string) string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprintf("%q", *v)
}
