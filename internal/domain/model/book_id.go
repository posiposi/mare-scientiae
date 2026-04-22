package model

import (
	"errors"
	"regexp"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

var (
	ErrBookIDRequired      = errors.New("book id: id is required")
	ErrBookIDInvalidFormat = errors.New("book id: id must be a valid UUID")
)

type BookID struct {
	value string
}

func NewBookID(v string) (BookID, error) {
	if v == "" {
		return BookID{}, ErrBookIDRequired
	}
	if !uuidRegex.MatchString(v) {
		return BookID{}, ErrBookIDInvalidFormat
	}
	return BookID{value: v}, nil
}

func (id BookID) String() string {
	return id.value
}
