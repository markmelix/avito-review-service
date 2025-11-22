package postgres

import (
	"context"
	"fmt"
	"review/internal/domain"
	"review/internal/repo"
)

func (pg *postgres) AddTeam(ctx context.Context, name string, members []domain.User) error {
	tx, err := pg.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error beginning postgres transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	res, err := tx.Exec(ctx, "INSERT INTO teams (name) VALUES ($1);", name)
	if err != nil {
		return fmt.Errorf("team insertion query: %w", err)
	}

	if res.RowsAffected() == 0 {
		return repo.ErrAlreadyExists
	}

	for _, u := range members {
		_, err := tx.Exec(
			ctx,
			"INSERT INTO users (id, username, is_active, team_name) VALUES ($1, $2, $3, $4);",
			u.Id, u.Username, u.IsActive, name,
		)
		if err != nil {
			return fmt.Errorf("user insertion query: %w", err)
		}
	}

	return err
}

func (pg *postgres) GetTeam(ctx context.Context, name string) (*domain.Team, error) {
	tx, err := pg.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error beginning postgres transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	var exists bool
	if err := tx.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM teams WHERE name = $1)", name).Scan(&exists); err != nil {
		return nil, fmt.Errorf("checking whether team exists: %w", err)
	}
	if !exists {
		return nil, repo.ErrNotFound
	}

	team := &domain.Team{Name: name}

	rows, err := tx.Query(ctx, "SELECT id, username, is_active FROM users WHERE team_name = $1", name)
	if err != nil {
		return nil, fmt.Errorf("while querying team members: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.Id, &u.Username, &u.IsActive); err != nil {
			return nil, fmt.Errorf("while scanning user: %w", err)
		}
		team.Members = append(team.Members, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("while iterating on the rows: %w", err)
	}

	return team, nil
}
