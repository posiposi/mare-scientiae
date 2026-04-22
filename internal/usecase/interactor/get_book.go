package interactor

import (
	"context"

	"helloworld/internal/domain/model"
	"helloworld/internal/domain/repository"
)

type GetBookInteractor struct {
	repo repository.BookQueryRepositorier
}

func NewGetBookInteractor(repo repository.BookQueryRepositorier) *GetBookInteractor {
	return &GetBookInteractor{repo: repo}
}

func (u *GetBookInteractor) Execute(ctx context.Context, id model.BookID) (*model.Book, error) {
	return u.repo.FindByID(ctx, id)
}
