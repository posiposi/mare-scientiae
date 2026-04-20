package usecase

import (
	"context"

	"helloworld/internal/domain"
)

type ListBooksUsecase struct {
	repo domain.BookQueryRepositorier
}

func NewListBooksUsecase(repo domain.BookQueryRepositorier) *ListBooksUsecase {
	return &ListBooksUsecase{repo: repo}
}

func (u *ListBooksUsecase) Execute(ctx context.Context) ([]*domain.Book, error) {
	return u.repo.FindAll(ctx)
}
