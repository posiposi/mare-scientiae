package domain

import "context"

// BookQueryRepositorier defines read-side persistence operations for
// Book aggregates (CQRS query side).
type BookQueryRepositorier interface {
	FindAll(ctx context.Context) ([]*Book, error)
}
