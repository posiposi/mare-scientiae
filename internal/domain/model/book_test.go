package model

import (
	"testing"
	"time"
)

func newTestBookID(t *testing.T, v string) BookID {
	t.Helper()
	id, err := NewBookID(v)
	if err != nil {
		t.Fatalf("NewBookID(%q) unexpected error: %v", v, err)
	}
	return id
}

func newTestGoogleBooksID(t *testing.T, v string) GoogleBooksID {
	t.Helper()
	gid, err := NewGoogleBooksID(v)
	if err != nil {
		t.Fatalf("NewGoogleBooksID(%q) unexpected error: %v", v, err)
	}
	return gid
}

func newTestBookTitle(t *testing.T, v string) BookTitle {
	t.Helper()
	title, err := NewBookTitle(v)
	if err != nil {
		t.Fatalf("NewBookTitle(%q) unexpected error: %v", v, err)
	}
	return title
}

func newTestBookSubtitle(t *testing.T, v string) *BookSubtitle {
	t.Helper()
	subtitle, err := NewBookSubtitle(&v)
	if err != nil {
		t.Fatalf("NewBookSubtitle(%q) unexpected error: %v", v, err)
	}
	return subtitle
}

func newTestAuthors(t *testing.T, names ...string) Authors {
	t.Helper()
	authors, err := NewAuthors(names)
	if err != nil {
		t.Fatalf("NewAuthors() unexpected error: %v", err)
	}
	return authors
}

func TestNewBook_正常系_全フィールド指定(t *testing.T) {
	now := time.Now()
	id := newTestBookID(t, "550e8400-e29b-41d4-a716-446655440000")
	gid := newTestGoogleBooksID(t, "googleId123")
	title := newTestBookTitle(t, "テストタイトル")
	subtitle := newTestBookSubtitle(t, "副題テスト")
	authors := newTestAuthors(t, "著者A", "著者B")

	book := NewBook(id, gid, title, subtitle, authors, now, now)

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
	if book.Authors.Values()[1].String() != "著者B" {
		t.Errorf("NewBook().Authors[1].String() = %q, want %q", book.Authors.Values()[1].String(), "著者B")
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
	id := newTestBookID(t, "550e8400-e29b-41d4-a716-446655440000")
	gid := newTestGoogleBooksID(t, "googleId123")
	title := newTestBookTitle(t, "テストタイトル")
	authors := newTestAuthors(t, "著者A")

	book := NewBook(id, gid, title, nil, authors, now, now)

	if book.Subtitle != nil {
		t.Errorf("NewBook().Subtitle = %v, want nil", book.Subtitle)
	}
}
