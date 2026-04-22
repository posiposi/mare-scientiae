package model

import "errors"

var (
	ErrBookTitleRequired = errors.New("book title: title is required")
	ErrBookTitleTooLong  = errors.New("book title: title must be 500 characters or less")
)

const bookTitleMaxLen = 500

type BookTitle struct {
	value string
}

func NewBookTitle(v string) (BookTitle, error) {
	if v == "" {
		return BookTitle{}, ErrBookTitleRequired
	}
	if len([]rune(v)) > bookTitleMaxLen {
		return BookTitle{}, ErrBookTitleTooLong
	}
	return BookTitle{value: v}, nil
}

func (t BookTitle) String() string {
	return t.value
}
