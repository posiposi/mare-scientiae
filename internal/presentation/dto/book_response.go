package dto

import (
	"time"

	"helloworld/internal/domain/model"
)

type BookResponse struct {
	ID            string    `json:"id"`
	GoogleBooksID string    `json:"google_books_id"`
	Title         string    `json:"title"`
	Subtitle      *string   `json:"subtitle"`
	Authors       []string  `json:"authors"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ListBooksResponse struct {
	Books []BookResponse `json:"books"`
}

func NewBookResponse(b *model.Book) BookResponse {
	var subtitle *string
	if b.Subtitle != nil {
		s := b.Subtitle.String()
		subtitle = &s
	}
	authorVOs := b.Authors.Values()
	authors := make([]string, 0, len(authorVOs))
	for _, a := range authorVOs {
		authors = append(authors, a.String())
	}
	return BookResponse{
		ID:            b.ID.String(),
		GoogleBooksID: b.GoogleBooksID.String(),
		Title:         b.Title.String(),
		Subtitle:      subtitle,
		Authors:       authors,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
	}
}

func NewListBooksResponse(books []*model.Book) ListBooksResponse {
	out := make([]BookResponse, 0, len(books))
	for _, b := range books {
		out = append(out, NewBookResponse(b))
	}
	return ListBooksResponse{Books: out}
}
