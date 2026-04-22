package repository

import (
	"context"

	"helloworld/internal/domain/model"
)

type BookQueryRepositorier interface {
	FindAll(ctx context.Context) ([]*model.Book, error)
}
