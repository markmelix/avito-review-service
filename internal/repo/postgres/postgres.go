package postgres

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dsn = fmt.Sprintf("postgresql://postgres:%s@postgres:5432/postgres", os.Getenv("POSTGRES_PASSWORD"))

func getConfig() *pgxpool.Config {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse db config: %v\n", err)
		os.Exit(1)
	}
	return config
}

type postgres struct {
	db *pgxpool.Pool
}

var (
	pgInstance *postgres
	pgOnce     sync.Once
)

func NewPostgres(ctx context.Context) *postgres {
	pgOnce.Do(func() {
		db, err := pgxpool.NewWithConfig(ctx, getConfig())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
			os.Exit(1)
		}

		pgInstance = &postgres{db}
	})

	return pgInstance
}

func (pg *postgres) Close() {
	pg.db.Close()
}

//go:embed queries/create_tables.sql
var createTablesQuery string

func (pg *postgres) CreateTables(ctx context.Context) error {
	_, err := pg.db.Exec(ctx, createTablesQuery)
	if err != nil {
		return fmt.Errorf("error processing create_tables query: %w", err)
	}
	return nil
}
