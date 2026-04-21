package dto_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"helloworld/internal/domain/model"
	"helloworld/internal/presentation/dto"
)

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

func TestNewListBooksResponse_ConvertsAllFields(t *testing.T) {
	createdAt := time.Date(2026, 4, 20, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 4, 21, 11, 0, 0, 0, time.UTC)
	subtitle := "Tackling Complexity"
	books := []*model.Book{
		buildBook(t, "11111111-1111-4111-8111-111111111111", "gbid-001", "DDD", &subtitle, []string{"Eric Evans"}, createdAt, updatedAt),
		buildBook(t, "22222222-2222-4222-8222-222222222222", "gbid-002", "TGPL", nil, []string{"Donovan", "Kernighan"}, createdAt, updatedAt),
	}

	got := dto.NewListBooksResponse(books)

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
			{
				ID:            "22222222-2222-4222-8222-222222222222",
				GoogleBooksID: "gbid-002",
				Title:         "TGPL",
				Subtitle:      nil,
				Authors:       []string{"Donovan", "Kernighan"},
				CreatedAt:     createdAt,
				UpdatedAt:     updatedAt,
			},
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("NewListBooksResponse() mismatch (-want +got):\n%s", diff)
	}
}

func TestNewListBooksResponse_EmptyInputReturnsEmptyBooks(t *testing.T) {
	got := dto.NewListBooksResponse([]*model.Book{})
	if got.Books == nil {
		t.Fatal("Books is nil, want non-nil empty slice")
	}
	if len(got.Books) != 0 {
		t.Errorf("len(Books) = %d, want 0", len(got.Books))
	}
}
