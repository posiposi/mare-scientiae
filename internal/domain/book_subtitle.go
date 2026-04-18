package domain

import "errors"

var (
	ErrBookSubtitleEmpty   = errors.New("book subtitle: subtitle must not be empty, use nil instead")
	ErrBookSubtitleTooLong = errors.New("book subtitle: subtitle must be 500 characters or less")
)

const bookSubtitleMaxLen = 500

type BookSubtitle struct {
	value string
}

func NewBookSubtitle(v *string) (*BookSubtitle, error) {
	if v == nil {
		return nil, nil
	}
	if *v == "" {
		return nil, ErrBookSubtitleEmpty
	}
	if len([]rune(*v)) > bookSubtitleMaxLen {
		return nil, ErrBookSubtitleTooLong
	}
	return &BookSubtitle{value: *v}, nil
}

func (s BookSubtitle) String() string {
	return s.value
}
