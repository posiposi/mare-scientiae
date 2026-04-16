package domain

import "errors"

var ErrAuthorsRequired = errors.New("authors: at least one author is required")

type Authors struct {
	values []Author
}

func NewAuthors(v []Author) (Authors, error) {
	if len(v) == 0 {
		return Authors{}, ErrAuthorsRequired
	}
	copied := make([]Author, len(v))
	copy(copied, v)
	return Authors{values: copied}, nil
}

func (a Authors) Values() []Author {
	copied := make([]Author, len(a.values))
	copy(copied, a.values)
	return copied
}
