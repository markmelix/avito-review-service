package postgres

import (
	"context"
	"fmt"
	"review/internal/domain"
	"review/internal/repo"

	"github.com/jackc/pgx/v5"
)

func userExists(ctx context.Context, tx pgx.Tx, id string) (bool, error) {
	var exists bool
	if err := tx.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM users WHERE id = $1);", id).Scan(&exists); err != nil {
		return false, fmt.Errorf("checking whether user exists: %w", err)
	}
	return exists, nil
}

func (pg *postgres) GetUser(ctx context.Context, id string) (*domain.User, error) {
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

	exists, err := userExists(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, repo.ErrNotFound
	}

	u := &domain.User{Id: id}

	tx.QueryRow(
		ctx,
		`
			SELECT u.username, u.is_active, t.name AS team_name
			FROM users u
			JOIN teams t ON t.id = u.team_id
			WHERE u.id = $1
		`,
		id,
	).Scan(&u.Username, &u.IsActive, &u.TeamName)
	if err != nil {
		return nil, fmt.Errorf("user is_active update query: %w", err)
	}

	return u, nil
}

func (pg *postgres) SetIsActiveUser(ctx context.Context, id string, isActive bool) error {
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

	exists, err := userExists(ctx, tx, id)
	if err != nil {
		return err
	}
	if !exists {
		return repo.ErrNotFound
	}

	_, err = tx.Exec(ctx, "UPDATE users SET is_active = $1 WHERE id = $2;", isActive, id)
	if err != nil {
		return fmt.Errorf("user is_active update query: %w", err)
	}

	return nil
}

func (pg *postgres) GetReviewAssignments(ctx context.Context, userId string) ([]domain.PullReq, error) {
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

	exists, err := userExists(ctx, tx, userId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, repo.ErrNotFound
	}

	prs := []domain.PullReq{}

	rows, err := tx.Query(ctx, "SELECT id, name, author_id, status FROM pull_requests WHERE author_id = $1;", userId)
	if err != nil {
		return nil, fmt.Errorf("while querying team members: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var pr domain.PullReq
		if err := rows.Scan(&pr.Id, &pr.Name, &pr.AuthorId, &pr.Status); err != nil {
			return nil, fmt.Errorf("while scanning pr: %w", err)
		}
		prs = append(prs, pr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("while iterating on the rows: %w", err)
	}

	return prs, nil
}
