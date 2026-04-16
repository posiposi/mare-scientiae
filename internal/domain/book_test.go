package domain

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewBook_正常系_全フィールド指定(t *testing.T) {
	subtitle := "副題テスト"
	now := time.Now()

	book, err := NewBook(
		"550e8400-e29b-41d4-a716-446655440000",
		"googleId123",
		"テストタイトル",
		&subtitle,
		[]string{"著者A", "著者B"},
		now,
		now,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if book.ID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("expected ID %q, got %q", "550e8400-e29b-41d4-a716-446655440000", book.ID)
	}
	if book.GoogleBooksID != "googleId123" {
		t.Errorf("expected GoogleBooksID %q, got %q", "googleId123", book.GoogleBooksID)
	}
	if book.Title != "テストタイトル" {
		t.Errorf("expected Title %q, got %q", "テストタイトル", book.Title)
	}
	if book.Subtitle == nil || *book.Subtitle != "副題テスト" {
		t.Errorf("expected Subtitle %q, got %v", "副題テスト", book.Subtitle)
	}
	if len(book.Authors) != 2 || book.Authors[0] != "著者A" || book.Authors[1] != "著者B" {
		t.Errorf("expected Authors [著者A 著者B], got %v", book.Authors)
	}
	if !book.CreatedAt.Equal(now) {
		t.Errorf("expected CreatedAt %v, got %v", now, book.CreatedAt)
	}
	if !book.UpdatedAt.Equal(now) {
		t.Errorf("expected UpdatedAt %v, got %v", now, book.UpdatedAt)
	}
}

func TestNewBook_正常系_Subtitleがnil(t *testing.T) {
	now := time.Now()

	book, err := NewBook(
		"550e8400-e29b-41d4-a716-446655440000",
		"googleId123",
		"テストタイトル",
		nil,
		[]string{"著者A"},
		now,
		now,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if book.Subtitle != nil {
		t.Errorf("expected Subtitle nil, got %v", book.Subtitle)
	}
}

func TestNewBook_異常系_IDが空文字(t *testing.T) {
	now := time.Now()

	_, err := NewBook("", "googleId123", "タイトル", nil, []string{"著者A"}, now, now)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrBookIDRequired) {
		t.Errorf("expected ErrBookIDRequired, got %v", err)
	}
}

func TestNewBook_異常系_IDが不正なUUID形式(t *testing.T) {
	now := time.Now()

	_, err := NewBook("not-a-uuid", "googleId123", "タイトル", nil, []string{"著者A"}, now, now)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrBookIDInvalidFormat) {
		t.Errorf("expected ErrBookIDInvalidFormat, got %v", err)
	}
}

func TestNewBook_異常系_GoogleBooksIDが空文字(t *testing.T) {
	now := time.Now()

	_, err := NewBook("550e8400-e29b-41d4-a716-446655440000", "", "タイトル", nil, []string{"著者A"}, now, now)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrBookGoogleBooksIDRequired) {
		t.Errorf("expected ErrBookGoogleBooksIDRequired, got %v", err)
	}
}

func TestNewBook_異常系_Titleが空文字(t *testing.T) {
	now := time.Now()

	_, err := NewBook("550e8400-e29b-41d4-a716-446655440000", "googleId123", "", nil, []string{"著者A"}, now, now)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrBookTitleRequired) {
		t.Errorf("expected ErrBookTitleRequired, got %v", err)
	}
}

func TestNewBook_異常系_Authorsが空またはnil(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		authors []string
	}{
		{"nil", nil},
		{"空スライス", []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewBook("550e8400-e29b-41d4-a716-446655440000", "googleId123", "タイトル", nil, tt.authors, now, now)

			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !errors.Is(err, ErrBookAuthorsRequired) {
				t.Errorf("expected ErrBookAuthorsRequired, got %v", err)
			}
		})
	}
}

func TestNewBook_異常系_Subtitleが空文字(t *testing.T) {
	now := time.Now()

	_, err := NewBook("550e8400-e29b-41d4-a716-446655440000", "googleId123", "タイトル", strPtr(""), []string{"著者A"}, now, now)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrBookSubtitleEmpty) {
		t.Errorf("expected ErrBookSubtitleEmpty, got %v", err)
	}
}

func TestNewBook_異常系_Authors内に空文字要素(t *testing.T) {
	now := time.Now()

	_, err := NewBook("550e8400-e29b-41d4-a716-446655440000", "googleId123", "タイトル", nil, []string{"著者A", ""}, now, now)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrBookAuthorEmpty) {
		t.Errorf("expected ErrBookAuthorEmpty, got %v", err)
	}
}

func TestNewBook_異常系_文字数上限超過(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		googleBooksID string
		title         string
		subtitle      *string
		expectedErr   error
	}{
		{
			name:          "GoogleBooksIDが51文字",
			googleBooksID: strings.Repeat("a", 51),
			title:         "タイトル",
			subtitle:      nil,
			expectedErr:   ErrBookGoogleBooksIDTooLong,
		},
		{
			name:          "Titleが501文字",
			googleBooksID: "googleId123",
			title:         strings.Repeat("あ", 501),
			subtitle:      nil,
			expectedErr:   ErrBookTitleTooLong,
		},
		{
			name:          "Subtitleが501文字",
			googleBooksID: "googleId123",
			title:         "タイトル",
			subtitle:      strPtr(strings.Repeat("あ", 501)),
			expectedErr:   ErrBookSubtitleTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewBook("550e8400-e29b-41d4-a716-446655440000", tt.googleBooksID, tt.title, tt.subtitle, []string{"著者A"}, now, now)

			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
