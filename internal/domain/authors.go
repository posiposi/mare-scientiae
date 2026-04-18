package domain

import "errors"

var ErrAuthorsRequired = errors.New("authors: at least one author is required")

type Authors struct {
	values []Author
}

func NewAuthors(values []string) (Authors, error) {
	if len(values) == 0 {
		return Authors{}, ErrAuthorsRequired
	}
	items := make([]Author, 0, len(values))
	for _, v := range values {
		a, err := NewAuthor(v)
		if err != nil {
			return Authors{}, err
		}
		items = append(items, a)
	}
	return Authors{values: items}, nil
}

func (a Authors) Values() []Author {
	copied := make([]Author, len(a.values))
	copy(copied, a.values)
	return copied
}
