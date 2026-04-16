package domain

import "errors"

var (
	ErrBookSubtitleEmpty   = errors.New("book subtitle: subtitle must not be empty, use nil instead")
	ErrBookSubtitleTooLong = errors.New("book subtitle: subtitle must be 500 characters or less")
)

type BookSubtitle struct {
	value string
}

func NewBookSubtitle(v string) (BookSubtitle, error) {
	if v == "" {
		return BookSubtitle{}, ErrBookSubtitleEmpty
	}
	if len([]rune(v)) > 500 {
		return BookSubtitle{}, ErrBookSubtitleTooLong
	}
	return BookSubtitle{value: v}, nil
}

func (s BookSubtitle) String() string {
	return s.value
}
