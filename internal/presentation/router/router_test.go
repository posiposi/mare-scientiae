package router_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"helloworld/internal/domain"
	"helloworld/internal/presentation/handler"
	"helloworld/internal/presentation/router"
)

type stubListBooksUsecase struct{}

func (stubListBooksUsecase) Execute(_ context.Context) ([]*domain.Book, error) {
	return []*domain.Book{}, nil
}

func newTestMux(t *testing.T) http.Handler {
	t.Helper()
	h := handler.NewBookHandler(stubListBooksUsecase{})
	return router.New(h)
}

func TestRouter_GetV1Books_Returns200(t *testing.T) {
	mux := newTestMux(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/books", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestRouter_PostV1Books_Returns405(t *testing.T) {
	mux := newTestMux(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/books", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestRouter_UnknownPath_Returns404(t *testing.T) {
	mux := newTestMux(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/unknown", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
