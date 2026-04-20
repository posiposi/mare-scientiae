package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"helloworld/internal/domain"
	"helloworld/internal/presentation/dto"
)

type ListBooksUsecaser interface {
	Execute(ctx context.Context) ([]*domain.Book, error)
}

type BookHandler struct {
	listBooksUsecase ListBooksUsecaser
}

func NewBookHandler(listBooksUsecase ListBooksUsecaser) *BookHandler {
	return &BookHandler{listBooksUsecase: listBooksUsecase}
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
