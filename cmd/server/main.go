package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"

	"helloworld/internal/infrastructure/ent"
	"helloworld/internal/infrastructure/persistence"
	"helloworld/internal/presentation/handler"
	"helloworld/internal/presentation/router"
	"helloworld/internal/usecase/interactor"
)

const listenAddr = ":8080"

var errDatabaseURLMissing = errors.New("DATABASE_URL is required")

func main() {
	if err := run(context.Background()); err != nil {
		log.Error().Err(err).Msg("server failed")
		os.Exit(1)
	}
}

func run(_ context.Context) error {
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

	bookRepo := persistence.NewBookRepository(client)
	listBooksInteractor := interactor.NewListBooksInteractor(bookRepo)
	bookHandler := handler.NewBookHandler(listBooksInteractor)
	mux := router.New(bookHandler)

	log.Info().Str("addr", listenAddr).Msg("starting http server")
	if err := http.ListenAndServe(listenAddr, mux); err != nil {
		return fmt.Errorf("listen and serve: %w", err)
	}
	return nil
}
