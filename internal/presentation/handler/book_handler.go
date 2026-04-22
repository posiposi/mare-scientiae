package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"

	"helloworld/internal/domain/model"
	"helloworld/internal/domain/repository"
	"helloworld/internal/presentation/dto"
)

type ListBooksUsecaser interface {
	Execute(ctx context.Context) ([]*model.Book, error)
}

type GetBookUsecaser interface {
	Execute(ctx context.Context, id model.BookID) (*model.Book, error)
}

type BookHandler struct {
	listBooksUsecase ListBooksUsecaser
	getBookUsecase   GetBookUsecaser
}

func NewBookHandler(listBooksUsecase ListBooksUsecaser, getBookUsecase GetBookUsecaser) *BookHandler {
	return &BookHandler{
		listBooksUsecase: listBooksUsecase,
		getBookUsecase:   getBookUsecase,
	}
}

func (h *BookHandler) ListBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.listBooksUsecase.Execute(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("list books usecase failed")
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal_server_error"})
		return
	}
	writeJSON(w, http.StatusOK, dto.NewListBooksResponse(books))
}

func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	id, err := model.NewBookID(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_book_id"})
		return
	}

	book, err := h.getBookUsecase.Execute(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrBookNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "book_not_found"})
			return
		}
		log.Error().Err(err).Str("book_id", id.String()).Msg("get book usecase failed")
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal_server_error"})
		return
	}
	writeJSON(w, http.StatusOK, dto.NewBookResponse(book))
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Error().Err(err).Msg("encode response")
	}
}
