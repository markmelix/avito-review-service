package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	handler "review/internal/http"
	"review/internal/repo/postgres"
)

var httpServerAddress string = "0.0.0.0:" + os.Getenv("SERVER_PORT")

func getPostgresPassword() string {
	p := os.Getenv("POSTGRES_PASSWORD")
	if p == "" {
		panic("POSTGRES_PASSWORD environment variable was not set")
	}
	return p
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx := context.Background()

	dsn := postgres.DsnFromPassword(getPostgresPassword())

	db := postgres.NewPostgres(ctx, dsn)
	if err := db.ApplyMigrations(ctx); err != nil {
		slog.Error("failed to apply database migrations", "error", err)
		return
	}
	defer db.Close()

	mux := handler.NewMux(db)

	slog.Info("Starting http-server", "address", httpServerAddress)
	if err := http.ListenAndServe(httpServerAddress, mux); err != nil {
		slog.Error("Error while serving http-server", "error", err)
	}
}
