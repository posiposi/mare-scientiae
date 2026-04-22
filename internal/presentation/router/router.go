package router

import (
	"net/http"

	"helloworld/internal/presentation/handler"
)

func New(bookHandler *handler.BookHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/books", bookHandler.ListBooks)
	return mux
}
