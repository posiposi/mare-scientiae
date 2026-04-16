package domain

import (
	"testing"
	"time"
)

func TestNewBook_正常系_全フィールド指定(t *testing.T) {
	now := time.Now()
	id, _ := NewBookID("550e8400-e29b-41d4-a716-446655440000")
	gid, _ := NewGoogleBooksID("googleId123")
	title, _ := NewBookTitle("テストタイトル")
	subtitle, _ := NewBookSubtitle("副題テスト")
	authorA, _ := NewAuthor("著者A")
	authorB, _ := NewAuthor("著者B")
	authors, _ := NewAuthors([]Author{authorA, authorB})

	book := NewBook(id, gid, title, &subtitle, authors, now, now)

	if book.ID.String() != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("NewBook().ID.String() = %q, want %q", book.ID.String(), "550e8400-e29b-41d4-a716-446655440000")
	}
	if book.GoogleBooksID.String() != "googleId123" {
		t.Errorf("NewBook().GoogleBooksID.String() = %q, want %q", book.GoogleBooksID.String(), "googleId123")
	}
	if book.Title.String() != "テストタイトル" {
		t.Errorf("NewBook().Title.String() = %q, want %q", book.Title.String(), "テストタイトル")
	}
	if book.Subtitle == nil || book.Subtitle.String() != "副題テスト" {
		t.Errorf("NewBook().Subtitle.String() = %v, want %q", book.Subtitle, "副題テスト")
	}
	if len(book.Authors.Values()) != 2 {
		t.Errorf("NewBook().Authors len = %d, want 2", len(book.Authors.Values()))
	}
	if book.Authors.Values()[0].String() != "著者A" {
		t.Errorf("NewBook().Authors[0].String() = %q, want %q", book.Authors.Values()[0].String(), "著者A")
	}
	if !book.CreatedAt.Equal(now) {
		t.Errorf("NewBook().CreatedAt = %v, want %v", book.CreatedAt, now)
	}
	if !book.UpdatedAt.Equal(now) {
		t.Errorf("NewBook().UpdatedAt = %v, want %v", book.UpdatedAt, now)
	}
}

func TestNewBook_正常系_Subtitleがnil(t *testing.T) {
	now := time.Now()
	id, _ := NewBookID("550e8400-e29b-41d4-a716-446655440000")
	gid, _ := NewGoogleBooksID("googleId123")
	title, _ := NewBookTitle("テストタイトル")
	authorA, _ := NewAuthor("著者A")
	authors, _ := NewAuthors([]Author{authorA})

	book := NewBook(id, gid, title, nil, authors, now, now)

	if book.Subtitle != nil {
		t.Errorf("NewBook().Subtitle = %v, want nil", book.Subtitle)
	}
}
