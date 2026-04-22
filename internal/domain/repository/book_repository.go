package repository

import (
	"context"
	"errors"

	"helloworld/internal/domain/model"
)

var ErrBookNotFound = errors.New("book: not found")

type BookQueryRepositorier interface {
	FindAll(ctx context.Context) ([]*model.Book, error)
	FindByID(ctx context.Context, id model.BookID) (*model.Book, error)
}
