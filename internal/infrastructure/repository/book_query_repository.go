package repository

import (
	"context"
	"fmt"

	"helloworld/internal/domain"
	"helloworld/internal/infrastructure/ent"
)

var _ domain.BookQueryRepositorier = (*BookQueryRepository)(nil)

type BookQueryRepository struct {
	client *ent.Client
}

func NewBookQueryRepository(client *ent.Client) *BookQueryRepository {
	return &BookQueryRepository{client: client}
}

func (r *BookQueryRepository) FindAll(ctx context.Context) ([]*domain.Book, error) {
	rows, err := r.client.Book.Query().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query all books: %w", err)
	}
	books := make([]*domain.Book, 0, len(rows))
	for _, row := range rows {
		b, err := toDomainBook(row)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func toDomainBook(row *ent.Book) (*domain.Book, error) {
	id, err := domain.NewBookID(row.ID.String())
	if err != nil {
		return nil, fmt.Errorf("book id (%s): %w", row.ID, err)
	}
	googleBooksID, err := domain.NewGoogleBooksID(row.GoogleBooksID)
	if err != nil {
		return nil, fmt.Errorf("google books id (%s): %w", row.GoogleBooksID, err)
	}
	title, err := domain.NewBookTitle(row.Title)
	if err != nil {
		return nil, fmt.Errorf("book title (%s): %w", row.ID, err)
	}
	subtitle, err := domain.NewBookSubtitle(row.Subtitle)
	if err != nil {
		return nil, fmt.Errorf("book subtitle (%s): %w", row.ID, err)
	}
	authors, err := domain.NewAuthors(row.Authors)
	if err != nil {
		return nil, fmt.Errorf("book authors (%s): %w", row.ID, err)
	}
	return domain.NewBook(id, googleBooksID, title, subtitle, authors, row.CreatedAt, row.UpdatedAt), nil
}
