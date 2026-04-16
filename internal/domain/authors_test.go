package domain

import (
	"errors"
	"testing"
)

func TestNewAuthors(t *testing.T) {
	authorA, _ := NewAuthor("著者A")
	authorB, _ := NewAuthor("著者B")

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
	authorA, _ := NewAuthor("著者A")
	authorB, _ := NewAuthor("著者B")
	authors, _ := NewAuthors([]Author{authorA, authorB})

	values := authors.Values()
	values[0] = Author{value: "改変"}

	if authors.Values()[0].String() != "著者A" {
		t.Errorf("Values() returned reference to internal slice, expected copy")
	}
}
