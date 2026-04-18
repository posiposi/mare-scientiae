package domain

import (
	"errors"
	"testing"
)

func TestNewAuthors(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    []string
		wantErr error
	}{
		{
			name:    "1人の著者",
			input:   []string{"著者A"},
			want:    []string{"著者A"},
			wantErr: nil,
		},
		{
			name:    "複数の著者",
			input:   []string{"著者A", "著者B"},
			want:    []string{"著者A", "著者B"},
			wantErr: nil,
		},
		{
			name:    "nilスライス",
			input:   nil,
			wantErr: ErrAuthorsRequired,
		},
		{
			name:    "空スライス",
			input:   []string{},
			wantErr: ErrAuthorsRequired,
		},
		{
			name:    "空文字を含む",
			input:   []string{"著者A", ""},
			wantErr: ErrAuthorEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAuthors(tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("NewAuthors(%v) error = %v, want %v", tt.input, err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("NewAuthors(%v) unexpected error: %v", tt.input, err)
			}
			values := got.Values()
			if len(values) != len(tt.want) {
				t.Fatalf("NewAuthors(%v) len = %d, want %d", tt.input, len(values), len(tt.want))
			}
			for i, v := range values {
				if v.String() != tt.want[i] {
					t.Errorf("NewAuthors(%v)[%d] = %q, want %q", tt.input, i, v.String(), tt.want[i])
				}
			}
		})
	}
}

func TestAuthors_Values_返却値の変更が元データに影響しない(t *testing.T) {
	authors, err := NewAuthors([]string{"著者A", "著者B"})
	if err != nil {
		t.Fatalf("NewAuthors() unexpected error: %v", err)
	}

	values := authors.Values()
	modified, err := NewAuthor("改変")
	if err != nil {
		t.Fatalf("NewAuthor(%q) unexpected error: %v", "改変", err)
	}
	values[0] = modified

	if authors.Values()[0].String() != "著者A" {
		t.Errorf("Authors.Values()[0].String() = %q, want %q", authors.Values()[0].String(), "著者A")
	}
}
