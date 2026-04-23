package persistence

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"helloworld/internal/domain/model"
	"helloworld/internal/domain/repository"
	"helloworld/internal/infrastructure/ent"
	"helloworld/internal/infrastructure/ent/book"
)

type BookRepository struct {
	client *ent.Client
}

func NewBookRepository(client *ent.Client) *BookRepository {
	return &BookRepository{client: client}
}

func (r *BookRepository) FindAll(ctx context.Context) ([]*model.Book, error) {
	rows, err := r.client.Book.Query().WithAuthors().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query all books: %w", err)
	}
	books := make([]*model.Book, 0, len(rows))
	for _, row := range rows {
		b, err := toDomainBook(row)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func (r *BookRepository) FindByID(ctx context.Context, id model.BookID) (*model.Book, error) {
	row, err := r.client.Book.Query().
		Where(book.IDEQ(uuid.MustParse(id.String()))).
		WithAuthors().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("find book (id=%s): %w", id.String(), repository.ErrBookNotFound)
		}
		return nil, fmt.Errorf("query book (id=%s): %w", id.String(), err)
	}
	return toDomainBook(row)
}

func toDomainBook(row *ent.Book) (*model.Book, error) {
	id, err := model.NewBookID(row.ID.String())
	if err != nil {
		return nil, fmt.Errorf("book id (value=%q): %w", row.ID, err)
	}
	googleBooksID, err := model.NewGoogleBooksID(row.GoogleBooksID)
	if err != nil {
		return nil, fmt.Errorf("google books id (id=%s, value=%q): %w", row.ID, row.GoogleBooksID, err)
	}
	title, err := model.NewBookTitle(row.Title)
	if err != nil {
		return nil, fmt.Errorf("book title (id=%s, value=%q): %w", row.ID, row.Title, err)
	}
	subtitle, err := model.NewBookSubtitle(row.Subtitle)
	if err != nil {
		return nil, fmt.Errorf("book subtitle (id=%s, value=%s): %w", row.ID, formatNillableString(row.Subtitle), err)
	}
	names := make([]string, 0, len(row.Edges.Authors))
	for _, a := range row.Edges.Authors {
		names = append(names, a.Name)
	}
	authors, err := model.NewAuthors(names)
	if err != nil {
		return nil, fmt.Errorf("book authors (id=%s, value=%v): %w", row.ID, names, err)
	}
	return model.NewBook(id, googleBooksID, title, subtitle, authors, row.CreatedAt, row.UpdatedAt), nil
}

func formatNillableString(v *string) string {
	if v == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%q", *v)
}
