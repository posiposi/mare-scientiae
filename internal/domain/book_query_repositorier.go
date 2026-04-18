package domain

import "context"

type BookQueryRepositorier interface {
	FindAll(ctx context.Context) ([]*Book, error)
}
