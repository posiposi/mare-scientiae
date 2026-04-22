package interactor_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"helloworld/internal/domain/model"
	"helloworld/internal/usecase/interactor"
)

type fakeBookQueryRepository struct {
	books []*model.Book
	err   error
}

func (f *fakeBookQueryRepository) FindAll(_ context.Context) ([]*model.Book, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.books, nil
}

func (f *fakeBookQueryRepository) FindByID(_ context.Context, _ model.BookID) (*model.Book, error) {
	return nil, nil
}

func newSampleBook(t *testing.T) *model.Book {
	t.Helper()
	id, err := model.NewBookID("11111111-1111-4111-8111-111111111111")
	if err != nil {
		t.Fatalf("NewBookID: %v", err)
	}
	gbid, err := model.NewGoogleBooksID("gbid-001")
	if err != nil {
		t.Fatalf("NewGoogleBooksID: %v", err)
	}
	title, err := model.NewBookTitle("Domain-Driven Design")
	if err != nil {
		t.Fatalf("NewBookTitle: %v", err)
	}
	subtitleStr := "Tackling Complexity in the Heart of Software"
	subtitle, err := model.NewBookSubtitle(&subtitleStr)
	if err != nil {
		t.Fatalf("NewBookSubtitle: %v", err)
	}
	authors, err := model.NewAuthors([]string{"Eric Evans"})
	if err != nil {
		t.Fatalf("NewAuthors: %v", err)
	}
	now := time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC)
	return model.NewBook(id, gbid, title, subtitle, authors, now, now)
}

func TestListBooksInteractor_Execute_ReturnsBooksFromRepository(t *testing.T) {
	book := newSampleBook(t)
	repo := &fakeBookQueryRepository{books: []*model.Book{book}}
	uc := interactor.NewListBooksInteractor(repo)

	got, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	want := []*model.Book{book}
	opts := cmp.AllowUnexported(
		model.Book{},
		model.BookID{},
		model.GoogleBooksID{},
		model.BookTitle{},
		model.BookSubtitle{},
		model.Authors{},
		model.Author{},
	)
	if diff := cmp.Diff(want, got, opts); diff != "" {
		t.Errorf("Execute() mismatch (-want +got):\n%s", diff)
	}
}

func TestListBooksInteractor_Execute_ReturnsEmptyWhenNoBooks(t *testing.T) {
	repo := &fakeBookQueryRepository{books: []*model.Book{}}
	uc := interactor.NewListBooksInteractor(repo)

	got, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if len(got) != 0 {
		t.Errorf("Execute() len = %d, want 0", len(got))
	}
}

func TestListBooksInteractor_Execute_PropagatesRepositoryError(t *testing.T) {
	wantErr := errors.New("repo failure")
	repo := &fakeBookQueryRepository{err: wantErr}
	uc := interactor.NewListBooksInteractor(repo)

	_, err := uc.Execute(context.Background())
	if !errors.Is(err, wantErr) {
		t.Errorf("Execute() error = %v, want %v", err, wantErr)
	}
}
