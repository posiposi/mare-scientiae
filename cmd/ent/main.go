package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"

	"helloworld/internal/infrastructure/ent"
)

var errDatabaseURLMissing = errors.New("DATABASE_URL is required")

func main() {
	if err := run(context.Background()); err != nil {
		log.Error().Err(err).Msg("migration failed")
		os.Exit(1)
	}
	log.Info().Msg("migration completed")
}

func run(ctx context.Context) error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return errDatabaseURLMissing
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	client := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, db)))
	defer client.Close()

	if err := client.Schema.Create(ctx, schema.WithDropColumn(true)); err != nil {
		return fmt.Errorf("schema create: %w", err)
	}
	return nil
}
