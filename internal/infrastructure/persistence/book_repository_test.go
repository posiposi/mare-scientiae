package persistence

import (
	"context"
	"errors"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"

	"helloworld/internal/domain/model"
	"helloworld/internal/domain/repository"
	"helloworld/internal/infrastructure/ent"
	"helloworld/internal/infrastructure/ent/author"
)

func TestBookRepository_FindAll_Empty(t *testing.T) {
	ctx := context.Background()
	truncateBooks(ctx, t)

	repo := NewBookRepository(testClient)
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

func TestBookRepository_FindAll_ReturnsInsertedBooks(t *testing.T) {
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
		authorEntities := ensureAuthors(ctx, t, s.authors)
		create := testClient.Book.Create().
			SetGoogleBooksID(s.googleBooksID).
			SetTitle(s.title).
			AddAuthors(authorEntities...)
		if s.subtitle != nil {
			create = create.SetSubtitle(*s.subtitle)
		}
		entBook, err := create.Save(ctx)
		if err != nil {
			t.Fatalf("seed insert (%s): %v", s.googleBooksID, err)
		}
		want[entBook.ID.String()] = s
	}

	repo := NewBookRepository(testClient)
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

func assertBookMatchesSeed(t *testing.T, got *model.Book, want map[string]bookSeed) {
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

func authorStrings(a model.Authors) []string {
	values := a.Values()
	out := make([]string, len(values))
	for i, v := range values {
		out[i] = v.String()
	}
	return out
}

func subtitleString(s *model.BookSubtitle) string {
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

func TestBookRepository_FindByID_ReturnsInsertedBook(t *testing.T) {
	ctx := context.Background()
	truncateBooks(ctx, t)

	subtitle := "Tackling Complexity in the Heart of Software"
	authorEntities := ensureAuthors(ctx, t, []string{"Eric Evans"})
	entBook, err := testClient.Book.Create().
		SetGoogleBooksID("gbid-ddd").
		SetTitle("Domain-Driven Design").
		SetSubtitle(subtitle).
		AddAuthors(authorEntities...).
		Save(ctx)
	if err != nil {
		t.Fatalf("seed insert: %v", err)
	}

	id, err := model.NewBookID(entBook.ID.String())
	if err != nil {
		t.Fatalf("NewBookID: %v", err)
	}

	repo := NewBookRepository(testClient)
	got, err := repo.FindByID(ctx, id)
	if err != nil {
		t.Fatalf("FindByID(%s) unexpected error: %v", id.String(), err)
	}

	if got.ID.String() != entBook.ID.String() {
		t.Errorf("FindByID().ID = %q, want %q", got.ID.String(), entBook.ID.String())
	}
	if got.GoogleBooksID.String() != "gbid-ddd" {
		t.Errorf("FindByID().GoogleBooksID = %q, want %q", got.GoogleBooksID.String(), "gbid-ddd")
	}
	if got.Title.String() != "Domain-Driven Design" {
		t.Errorf("FindByID().Title = %q, want %q", got.Title.String(), "Domain-Driven Design")
	}
	if g := subtitleString(got.Subtitle); g != subtitle {
		t.Errorf("FindByID().Subtitle = %q, want %q", g, subtitle)
	}
	if diff := cmp.Diff([]string{"Eric Evans"}, authorStrings(got.Authors)); diff != "" {
		t.Errorf("FindByID().Authors mismatch (-want +got):\n%s", diff)
	}
}

func TestBookRepository_FindByID_ReturnsErrBookNotFoundWhenMissing(t *testing.T) {
	ctx := context.Background()
	truncateBooks(ctx, t)

	id, err := model.NewBookID("22222222-2222-4222-8222-222222222222")
	if err != nil {
		t.Fatalf("NewBookID: %v", err)
	}

	repo := NewBookRepository(testClient)
	_, err = repo.FindByID(ctx, id)
	if !errors.Is(err, repository.ErrBookNotFound) {
		t.Errorf("FindByID() error = %v, want chain with %v", err, repository.ErrBookNotFound)
	}
}

func TestBookRepository_FindByID_ErrorOnInvalidPersistedData(t *testing.T) {
	ctx := context.Background()
	truncateBooks(ctx, t)

	authorEntities := ensureAuthors(ctx, t, []string{"Author"})
	entBook, err := testClient.Book.Create().
		SetGoogleBooksID("gbid-empty-subtitle").
		SetTitle("Title").
		SetSubtitle("").
		AddAuthors(authorEntities...).
		Save(ctx)
	if err != nil {
		t.Fatalf("seed insert: %v", err)
	}

	id, err := model.NewBookID(entBook.ID.String())
	if err != nil {
		t.Fatalf("NewBookID: %v", err)
	}

	repo := NewBookRepository(testClient)
	_, err = repo.FindByID(ctx, id)
	if !errors.Is(err, model.ErrBookSubtitleEmpty) {
		t.Errorf("FindByID() error = %v, want chain with %v", err, model.ErrBookSubtitleEmpty)
	}
}

func TestBookRepository_FindAll_ErrorOnInvalidPersistedData(t *testing.T) {
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
					Save(ctx); err != nil {
					t.Fatalf("seed insert: %v", err)
				}
			},
			wantErr: model.ErrAuthorsRequired,
		},
		{
			name: "authorsに空文字を含むレコード",
			seed: func(t *testing.T) {
				t.Helper()
				authorEntities := ensureAuthors(ctx, t, []string{"valid", ""})
				if _, err := testClient.Book.Create().
					SetGoogleBooksID("gbid-empty-author-element").
					SetTitle("Title").
					AddAuthors(authorEntities...).
					Save(ctx); err != nil {
					t.Fatalf("seed insert: %v", err)
				}
			},
			wantErr: model.ErrAuthorEmpty,
		},
		{
			name: "subtitleが空文字のレコード",
			seed: func(t *testing.T) {
				t.Helper()
				authorEntities := ensureAuthors(ctx, t, []string{"Author"})
				if _, err := testClient.Book.Create().
					SetGoogleBooksID("gbid-empty-subtitle").
					SetTitle("Title").
					SetSubtitle("").
					AddAuthors(authorEntities...).
					Save(ctx); err != nil {
					t.Fatalf("seed insert: %v", err)
				}
			},
			wantErr: model.ErrBookSubtitleEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			truncateBooks(ctx, t)
			tt.seed(t)

			repo := NewBookRepository(testClient)
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

func TestBookRepository_FindAll_SharesAuthorMaster(t *testing.T) {
	ctx := context.Background()
	truncateBooks(ctx, t)

	sharedAuthor := "Eric Evans"
	seeds := []bookSeed{
		{
			googleBooksID: "gbid-shared-1",
			title:         "Book One",
			authors:       []string{sharedAuthor},
		},
		{
			googleBooksID: "gbid-shared-2",
			title:         "Book Two",
			authors:       []string{sharedAuthor},
		},
	}
	for _, s := range seeds {
		authorEntities := ensureAuthors(ctx, t, s.authors)
		if _, err := testClient.Book.Create().
			SetGoogleBooksID(s.googleBooksID).
			SetTitle(s.title).
			AddAuthors(authorEntities...).
			Save(ctx); err != nil {
			t.Fatalf("seed insert (%s): %v", s.googleBooksID, err)
		}
	}

	repo := NewBookRepository(testClient)
	got, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("FindAll() unexpected error: %v", err)
	}
	if len(got) != len(seeds) {
		t.Fatalf("FindAll() len = %d, want %d", len(got), len(seeds))
	}
	for _, b := range got {
		gotAuthors := authorStrings(b.Authors)
		if diff := cmp.Diff([]string{sharedAuthor}, gotAuthors); diff != "" {
			t.Errorf("book(%s).Authors mismatch (-want +got):\n%s", b.ID.String(), diff)
		}
	}

	count, err := testClient.Author.Query().Count(ctx)
	if err != nil {
		t.Fatalf("count authors: %v", err)
	}
	if count != 1 {
		t.Errorf("Author count = %d, want 1 (master should be shared across books)", count)
	}
}

func ensureAuthors(ctx context.Context, t *testing.T, names []string) []*ent.Author {
	t.Helper()
	authors := make([]*ent.Author, 0, len(names))
	for _, name := range names {
		a, err := testClient.Author.Query().
			Where(author.NameEQ(name)).
			Only(ctx)
		if ent.IsNotFound(err) {
			a, err = testClient.Author.Create().SetName(name).Save(ctx)
		}
		if err != nil {
			t.Fatalf("ensure author %q: %v", name, err)
		}
		authors = append(authors, a)
	}
	return authors
}
