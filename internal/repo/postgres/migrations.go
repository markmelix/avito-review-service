package postgres

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	mpgx "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var migrations embed.FS

func (pg *postgres) ApplyMigrations(ctx context.Context) error {
	fsDriver, err := iofs.New(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("getting migrations embedding: %w", err)
	}

	db := sql.OpenDB(stdlib.GetPoolConnector(pg.db))

	dbDriver, err := mpgx.WithInstance(db, &mpgx.Config{})
	if err != nil {
		return fmt.Errorf("getting migrate instance from pgxpool: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", fsDriver, "postgres", dbDriver)
	if err != nil {
		return fmt.Errorf("getting migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("applying migrations: %w", err)
	}

	return nil
}
