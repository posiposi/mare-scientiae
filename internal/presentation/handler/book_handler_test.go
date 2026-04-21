package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"helloworld/internal/domain/model"
	"helloworld/internal/presentation/dto"
	"helloworld/internal/presentation/handler"
)

type fakeListBooksUsecase struct {
	books []*model.Book
	err   error
}

func (f *fakeListBooksUsecase) Execute(_ context.Context) ([]*model.Book, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.books, nil
}

func buildBook(t *testing.T, id, gbid, title string, subtitle *string, authorNames []string, createdAt, updatedAt time.Time) *model.Book {
	t.Helper()
	bookID, err := model.NewBookID(id)
	if err != nil {
		t.Fatalf("NewBookID: %v", err)
	}
	googleBooksID, err := model.NewGoogleBooksID(gbid)
	if err != nil {
		t.Fatalf("NewGoogleBooksID: %v", err)
	}
	bookTitle, err := model.NewBookTitle(title)
	if err != nil {
		t.Fatalf("NewBookTitle: %v", err)
	}
	bookSubtitle, err := model.NewBookSubtitle(subtitle)
	if err != nil {
		t.Fatalf("NewBookSubtitle: %v", err)
	}
	authors, err := model.NewAuthors(authorNames)
	if err != nil {
		t.Fatalf("NewAuthors: %v", err)
	}
	return model.NewBook(bookID, googleBooksID, bookTitle, bookSubtitle, authors, createdAt, updatedAt)
}

func TestBookHandler_ListBooks_Returns200WithBooksJSON(t *testing.T) {
	createdAt := time.Date(2026, 4, 20, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 4, 21, 11, 0, 0, 0, time.UTC)
	subtitle := "Tackling Complexity"
	book := buildBook(t, "11111111-1111-4111-8111-111111111111", "gbid-001", "DDD", &subtitle, []string{"Eric Evans"}, createdAt, updatedAt)

	uc := &fakeListBooksUsecase{books: []*model.Book{book}}
	h := handler.NewBookHandler(uc)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/books", nil)
	h.ListBooks(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}

	var got dto.ListBooksResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	want := dto.ListBooksResponse{
		Books: []dto.BookResponse{
			{
				ID:            "11111111-1111-4111-8111-111111111111",
				GoogleBooksID: "gbid-001",
				Title:         "DDD",
				Subtitle:      &subtitle,
				Authors:       []string{"Eric Evans"},
				CreatedAt:     createdAt,
				UpdatedAt:     updatedAt,
			},
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("body mismatch (-want +got):\n%s", diff)
	}
}

func TestBookHandler_ListBooks_ReturnsEmptyArrayWhenNoBooks(t *testing.T) {
	uc := &fakeListBooksUsecase{books: []*model.Book{}}
	h := handler.NewBookHandler(uc)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/books", nil)
	h.ListBooks(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Body.String(); got != `{"books":[]}`+"\n" && got != `{"books":[]}` {
		t.Errorf("body = %q, want empty books envelope", got)
	}
}

func TestBookHandler_ListBooks_Returns500WhenUsecaseFails(t *testing.T) {
	uc := &fakeListBooksUsecase{err: errors.New("boom")}
	h := handler.NewBookHandler(uc)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/books", nil)
	h.ListBooks(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}
}
