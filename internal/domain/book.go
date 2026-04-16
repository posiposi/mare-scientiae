package domain

import "time"

type Book struct {
	ID            BookID
	GoogleBooksID GoogleBooksID
	Title         BookTitle
	Subtitle      *BookSubtitle
	Authors       Authors
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewBook(id BookID, googleBooksID GoogleBooksID, title BookTitle, subtitle *BookSubtitle, authors Authors, createdAt, updatedAt time.Time) *Book {
	return &Book{
		ID:            id,
		GoogleBooksID: googleBooksID,
		Title:         title,
		Subtitle:      subtitle,
		Authors:       authors,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}
