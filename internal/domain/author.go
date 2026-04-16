package domain

import "errors"

var ErrAuthorEmpty = errors.New("author: name must not be empty")

type Author struct {
	value string
}

func NewAuthor(v string) (Author, error) {
	if v == "" {
		return Author{}, ErrAuthorEmpty
	}
	return Author{value: v}, nil
}

func (a Author) String() string {
	return a.value
}
