package interactor_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"

	"helloworld/internal/domain/model"
	"helloworld/internal/domain/repository"
	"helloworld/internal/usecase/interactor"
)

type fakeBookFinder struct {
	book *model.Book
	err  error
}

func (f *fakeBookFinder) FindAll(_ context.Context) ([]*model.Book, error) {
	return nil, nil
}

func (f *fakeBookFinder) FindByID(_ context.Context, _ model.BookID) (*model.Book, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.book, nil
}

func TestGetBookInteractor_Execute_ReturnsBookFromRepository(t *testing.T) {
	book := newSampleBook(t)
	repo := &fakeBookFinder{book: book}
	uc := interactor.NewGetBookInteractor(repo)

	got, err := uc.Execute(context.Background(), book.ID)
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	opts := cmp.AllowUnexported(
		model.Book{},
		model.BookID{},
		model.GoogleBooksID{},
		model.BookTitle{},
		model.BookSubtitle{},
		model.Authors{},
		model.Author{},
	)
	if diff := cmp.Diff(book, got, opts); diff != "" {
		t.Errorf("Execute() mismatch (-want +got):\n%s", diff)
	}
}

func TestGetBookInteractor_Execute_PropagatesNotFoundError(t *testing.T) {
	repo := &fakeBookFinder{err: repository.ErrBookNotFound}
	uc := interactor.NewGetBookInteractor(repo)

	id, err := model.NewBookID("33333333-3333-4333-8333-333333333333")
	if err != nil {
		t.Fatalf("NewBookID: %v", err)
	}

	_, err = uc.Execute(context.Background(), id)
	if !errors.Is(err, repository.ErrBookNotFound) {
		t.Errorf("Execute() error = %v, want chain with %v", err, repository.ErrBookNotFound)
	}
}

func TestGetBookInteractor_Execute_PropagatesRepositoryError(t *testing.T) {
	wantErr := errors.New("repo failure")
	repo := &fakeBookFinder{err: wantErr}
	uc := interactor.NewGetBookInteractor(repo)

	id, err := model.NewBookID("44444444-4444-4444-8444-444444444444")
	if err != nil {
		t.Fatalf("NewBookID: %v", err)
	}

	_, err = uc.Execute(context.Background(), id)
	if !errors.Is(err, wantErr) {
		t.Errorf("Execute() error = %v, want %v", err, wantErr)
	}
}
