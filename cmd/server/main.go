package main

import (
	"context"
	"log"
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
	log.SetFlags(log.LstdFlags)

	ctx := context.Background()

	dsn := postgres.DsnFromPassword(getPostgresPassword())

	db := postgres.NewPostgres(ctx, dsn)
	if err := db.CreateTables(ctx); err != nil {
		log.Fatalf("failed to create database tables: %v", err)
	}
	defer db.Close()

	mux := handler.NewMux(db)

	log.Printf("Starting http-server on %s\n", httpServerAddress)
	if err := http.ListenAndServe(httpServerAddress, mux); err != nil {
		log.Fatal("Error while serving http-server: %w", err)
	}
}
