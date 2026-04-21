package repository

import (
	"context"

	"helloworld/internal/domain/model"
)

// BookQueryRepositorier defines read-side persistence operations for
// Book aggregates (CQRS query side).
type BookQueryRepositorier interface {
	FindAll(ctx context.Context) ([]*model.Book, error)
}
