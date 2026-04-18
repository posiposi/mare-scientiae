package repository

import (
	"context"
	"errors"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"

	"helloworld/internal/domain"
)

func TestBookQueryRepository_FindAll_Empty(t *testing.T) {
	ctx := context.Background()
	truncateBooks(ctx, t)

	repo := NewBookQueryRepository(testClient)
	got, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("FindAll() unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("FindAll() len = %d, want 0", len(got))
	}
}

type bookSeed struct {
	googleBooksID string
	title         string
	subtitle      *string
	authors       []string
}

func TestBookQueryRepository_FindAll_ReturnsInsertedBooks(t *testing.T) {
	ctx := context.Background()
	truncateBooks(ctx, t)

	subtitleA := "Volume 1"
	seeds := []bookSeed{
		{
			googleBooksID: "gbid-001",
			title:         "Domain-Driven Design",
			subtitle:      &subtitleA,
			authors:       []string{"Eric Evans"},
		},
		{
			googleBooksID: "gbid-002",
			title:         "The Go Programming Language",
			subtitle:      nil,
			authors:       []string{"Alan A. A. Donovan", "Brian W. Kernighan"},
		},
	}

	want := make(map[string]bookSeed, len(seeds))
	for _, s := range seeds {
		create := testClient.Book.Create().
			SetGoogleBooksID(s.googleBooksID).
			SetTitle(s.title).
			SetAuthors(s.authors)
		if s.subtitle != nil {
			create = create.SetSubtitle(*s.subtitle)
		}
		entBook, err := create.Save(ctx)
		if err != nil {
			t.Fatalf("seed insert (%s): %v", s.googleBooksID, err)
		}
		want[entBook.ID.String()] = s
	}

	repo := NewBookQueryRepository(testClient)
	got, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("FindAll() unexpected error: %v", err)
	}
	if len(got) != len(seeds) {
		t.Fatalf("FindAll() len = %d, want %d", len(got), len(seeds))
	}

	for _, b := range got {
		assertBookMatchesSeed(t, b, want)
	}
}

func assertBookMatchesSeed(t *testing.T, got *domain.Book, want map[string]bookSeed) {
	t.Helper()

	id := got.ID.String()
	seed, ok := want[id]
	if !ok {
		t.Errorf("FindAll() returned unexpected book id=%s", id)
		return
	}

	if got.GoogleBooksID.String() != seed.googleBooksID {
		t.Errorf("book(%s).GoogleBooksID = %q, want %q", id, got.GoogleBooksID.String(), seed.googleBooksID)
	}
	if got.Title.String() != seed.title {
		t.Errorf("book(%s).Title = %q, want %q", id, got.Title.String(), seed.title)
	}
	if g, w := subtitleString(got.Subtitle), derefString(seed.subtitle); g != w {
		t.Errorf("book(%s).Subtitle = %q, want %q", id, g, w)
	}

	gotAuthors := authorStrings(got.Authors)
	wantAuthors := append([]string(nil), seed.authors...)
	sort.Strings(gotAuthors)
	sort.Strings(wantAuthors)
	if diff := cmp.Diff(wantAuthors, gotAuthors); diff != "" {
		t.Errorf("book(%s).Authors mismatch (-want +got):\n%s", id, diff)
	}
}

func authorStrings(a domain.Authors) []string {
	values := a.Values()
	out := make([]string, len(values))
	for i, v := range values {
		out[i] = v.String()
	}
	return out
}

func subtitleString(s *domain.BookSubtitle) string {
	if s == nil {
		return ""
	}
	return s.String()
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func TestBookQueryRepository_FindAll_ErrorOnInvalidPersistedData(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		seed    func(t *testing.T)
		wantErr error
	}{
		{
			name: "authorsが空配列のレコード",
			seed: func(t *testing.T) {
				t.Helper()
				if _, err := testClient.Book.Create().
					SetGoogleBooksID("gbid-empty-authors").
					SetTitle("Title").
					SetAuthors([]string{}).
					Save(ctx); err != nil {
					t.Fatalf("seed insert: %v", err)
				}
			},
			wantErr: domain.ErrAuthorsRequired,
		},
		{
			name: "authorsに空文字を含むレコード",
			seed: func(t *testing.T) {
				t.Helper()
				if _, err := testClient.Book.Create().
					SetGoogleBooksID("gbid-empty-author-element").
					SetTitle("Title").
					SetAuthors([]string{"valid", ""}).
					Save(ctx); err != nil {
					t.Fatalf("seed insert: %v", err)
				}
			},
			wantErr: domain.ErrAuthorEmpty,
		},
		{
			name: "subtitleが空文字のレコード",
			seed: func(t *testing.T) {
				t.Helper()
				if _, err := testClient.Book.Create().
					SetGoogleBooksID("gbid-empty-subtitle").
					SetTitle("Title").
					SetSubtitle("").
					SetAuthors([]string{"Author"}).
					Save(ctx); err != nil {
					t.Fatalf("seed insert: %v", err)
				}
			},
			wantErr: domain.ErrBookSubtitleEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			truncateBooks(ctx, t)
			tt.seed(t)

			repo := NewBookQueryRepository(testClient)
			_, err := repo.FindAll(ctx)
			if err == nil {
				t.Fatalf("FindAll() error = nil, want %v", tt.wantErr)
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("FindAll() error = %v, want chain with %v", err, tt.wantErr)
			}
		})
	}
}
