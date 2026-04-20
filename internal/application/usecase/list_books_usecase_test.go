package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"helloworld/internal/application/usecase"
	"helloworld/internal/domain"
)

type fakeBookQueryRepository struct {
	books []*domain.Book
	err   error
}

func (f *fakeBookQueryRepository) FindAll(_ context.Context) ([]*domain.Book, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.books, nil
}

func newSampleBook(t *testing.T) *domain.Book {
	t.Helper()
	id, err := domain.NewBookID("11111111-1111-4111-8111-111111111111")
	if err != nil {
		t.Fatalf("NewBookID: %v", err)
	}
	gbid, err := domain.NewGoogleBooksID("gbid-001")
	if err != nil {
		t.Fatalf("NewGoogleBooksID: %v", err)
	}
	title, err := domain.NewBookTitle("Domain-Driven Design")
	if err != nil {
		t.Fatalf("NewBookTitle: %v", err)
	}
	subtitleStr := "Tackling Complexity in the Heart of Software"
	subtitle, err := domain.NewBookSubtitle(&subtitleStr)
	if err != nil {
		t.Fatalf("NewBookSubtitle: %v", err)
	}
	authors, err := domain.NewAuthors([]string{"Eric Evans"})
	if err != nil {
		t.Fatalf("NewAuthors: %v", err)
	}
	now := time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC)
	return domain.NewBook(id, gbid, title, subtitle, authors, now, now)
}

func TestListBooksUsecase_Execute_ReturnsBooksFromRepository(t *testing.T) {
	book := newSampleBook(t)
	repo := &fakeBookQueryRepository{books: []*domain.Book{book}}
	uc := usecase.NewListBooksUsecase(repo)

	got, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	want := []*domain.Book{book}
	opts := cmp.AllowUnexported(
		domain.Book{},
		domain.BookID{},
		domain.GoogleBooksID{},
		domain.BookTitle{},
		domain.BookSubtitle{},
		domain.Authors{},
		domain.Author{},
	)
	if diff := cmp.Diff(want, got, opts); diff != "" {
		t.Errorf("Execute() mismatch (-want +got):\n%s", diff)
	}
}

func TestListBooksUsecase_Execute_ReturnsEmptyWhenNoBooks(t *testing.T) {
	repo := &fakeBookQueryRepository{books: []*domain.Book{}}
	uc := usecase.NewListBooksUsecase(repo)

	got, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if len(got) != 0 {
		t.Errorf("Execute() len = %d, want 0", len(got))
	}
}

func TestListBooksUsecase_Execute_PropagatesRepositoryError(t *testing.T) {
	wantErr := errors.New("repo failure")
	repo := &fakeBookQueryRepository{err: wantErr}
	uc := usecase.NewListBooksUsecase(repo)

	_, err := uc.Execute(context.Background())
	if !errors.Is(err, wantErr) {
		t.Errorf("Execute() error = %v, want %v", err, wantErr)
	}
}
