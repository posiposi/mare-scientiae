package domain

import (
	"errors"
	"testing"
)

func TestNewAuthors(t *testing.T) {
	authorA, err := NewAuthor("著者A")
	if err != nil {
		t.Fatalf("NewAuthor(%q) unexpected error: %v", "著者A", err)
	}
	authorB, err := NewAuthor("著者B")
	if err != nil {
		t.Fatalf("NewAuthor(%q) unexpected error: %v", "著者B", err)
	}

	tests := []struct {
		name    string
		input   []Author
		want    int
		wantErr error
	}{
		{
			name:    "1人の著者",
			input:   []Author{authorA},
			want:    1,
			wantErr: nil,
		},
		{
			name:    "複数の著者",
			input:   []Author{authorA, authorB},
			want:    2,
			wantErr: nil,
		},
		{
			name:    "nilスライス",
			input:   nil,
			want:    0,
			wantErr: ErrAuthorsRequired,
		},
		{
			name:    "空スライス",
			input:   []Author{},
			want:    0,
			wantErr: ErrAuthorsRequired,
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
					t.Errorf("NewAuthors() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("NewAuthors() unexpected error: %v", err)
			}
			if len(got.Values()) != tt.want {
				t.Errorf("NewAuthors() len = %d, want %d", len(got.Values()), tt.want)
			}
		})
	}
}

func TestAuthors_Values_返却値の変更が元データに影響しない(t *testing.T) {
	authorA, err := NewAuthor("著者A")
	if err != nil {
		t.Fatalf("NewAuthor(%q) unexpected error: %v", "著者A", err)
	}
	authorB, err := NewAuthor("著者B")
	if err != nil {
		t.Fatalf("NewAuthor(%q) unexpected error: %v", "著者B", err)
	}
	authors, err := NewAuthors([]Author{authorA, authorB})
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
