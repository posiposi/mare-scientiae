package domain

import (
	"errors"
	"regexp"
	"time"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

var (
	ErrBookIDRequired            = errors.New("book: id is required")
	ErrBookIDInvalidFormat       = errors.New("book: id must be a valid UUID")
	ErrBookGoogleBooksIDRequired = errors.New("book: google_books_id is required")
	ErrBookTitleRequired         = errors.New("book: title is required")
	ErrBookAuthorsRequired       = errors.New("book: authors is required")
	ErrBookAuthorEmpty           = errors.New("book: author must not be empty")
	ErrBookGoogleBooksIDTooLong  = errors.New("book: google_books_id must be 50 characters or less")
	ErrBookTitleTooLong          = errors.New("book: title must be 500 characters or less")
	ErrBookSubtitleTooLong       = errors.New("book: subtitle must be 500 characters or less")
	ErrBookSubtitleEmpty         = errors.New("book: subtitle must not be empty, use nil instead")
)

type Book struct {
	ID            string
	GoogleBooksID string
	Title         string
	Subtitle      *string
	Authors       []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewBook(id, googleBooksID, title string, subtitle *string, authors []string, createdAt, updatedAt time.Time) (*Book, error) {
	if id == "" {
		return nil, ErrBookIDRequired
	}
	if !uuidRegex.MatchString(id) {
		return nil, ErrBookIDInvalidFormat
	}
	if googleBooksID == "" {
		return nil, ErrBookGoogleBooksIDRequired
	}
	if len([]rune(googleBooksID)) > 50 {
		return nil, ErrBookGoogleBooksIDTooLong
	}
	if title == "" {
		return nil, ErrBookTitleRequired
	}
	if len([]rune(title)) > 500 {
		return nil, ErrBookTitleTooLong
	}
	if subtitle != nil {
		if *subtitle == "" {
			return nil, ErrBookSubtitleEmpty
		}
		if len([]rune(*subtitle)) > 500 {
			return nil, ErrBookSubtitleTooLong
		}
	}
	if len(authors) == 0 {
		return nil, ErrBookAuthorsRequired
	}
	for _, a := range authors {
		if a == "" {
			return nil, ErrBookAuthorEmpty
		}
	}

	return &Book{
		ID:            id,
		GoogleBooksID: googleBooksID,
		Title:         title,
		Subtitle:      subtitle,
		Authors:       authors,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}, nil
}
