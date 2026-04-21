package interactor

import (
	"context"

	"helloworld/internal/domain/model"
	"helloworld/internal/domain/repository"
)

type ListBooksInteractor struct {
	repo repository.BookQueryRepositorier
}

func NewListBooksInteractor(repo repository.BookQueryRepositorier) *ListBooksInteractor {
	return &ListBooksInteractor{repo: repo}
}

func (u *ListBooksInteractor) Execute(ctx context.Context) ([]*model.Book, error) {
	return u.repo.FindAll(ctx)
}
