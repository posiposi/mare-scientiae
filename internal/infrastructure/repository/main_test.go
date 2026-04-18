package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"

	"helloworld/internal/infrastructure/ent"
)

var (
	errTestDatabaseURLMissing = errors.New("TEST_DATABASE_URL is required")
	testClient                *ent.Client
)

func TestMain(m *testing.M) {
	code, err := runTestMain(m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "test setup: %v\n", err)
		os.Exit(1)
	}
	os.Exit(code)
}

func runTestMain(m *testing.M) (int, error) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		return 0, errTestDatabaseURLMissing
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return 0, fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	testClient = ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, db)))
	defer testClient.Close()

	if err := testClient.Schema.Create(context.Background()); err != nil {
		return 0, fmt.Errorf("schema create: %w", err)
	}

	return m.Run(), nil
}

func truncateBooks(ctx context.Context, t *testing.T) {
	t.Helper()
	if _, err := testClient.Book.Delete().Exec(ctx); err != nil {
		t.Fatalf("truncate books: %v", err)
	}
}
