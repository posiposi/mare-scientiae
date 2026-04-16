package domain

import "errors"

var (
	ErrGoogleBooksIDRequired = errors.New("google books id: id is required")
	ErrGoogleBooksIDTooLong  = errors.New("google books id: id must be 50 characters or less")
)

type GoogleBooksID struct {
	value string
}

func NewGoogleBooksID(v string) (GoogleBooksID, error) {
	if v == "" {
		return GoogleBooksID{}, ErrGoogleBooksIDRequired
	}
	if len([]rune(v)) > 50 {
		return GoogleBooksID{}, ErrGoogleBooksIDTooLong
	}
	return GoogleBooksID{value: v}, nil
}

func (g GoogleBooksID) String() string {
	return g.value
}
