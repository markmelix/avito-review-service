package postgres

import (
	"context"
	"fmt"
	"review/internal/domain"
	"review/internal/repo"

	"github.com/jackc/pgx/v5"
)

func prExists(ctx context.Context, tx pgx.Tx, id string) (bool, error) {
	var exists bool
	if err := tx.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM pull_requests WHERE id = $1);", id).Scan(&exists); err != nil {
		return false, fmt.Errorf("checking whether pr exists: %w", err)
	}
	return exists, nil
}

func (pg *postgres) CreatePullReq(ctx context.Context, pullReqId, name, authorId string) error {
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

	res, err := tx.Exec(
		ctx,
		"INSERT INTO pull_requests (id, name, author_id, status) VALUES ($1, $2, $3, $4);",
		pullReqId, name, authorId, domain.PullReqOpen,
	)
	if err != nil {
		return fmt.Errorf("pr insertion query: %w", err)
	}

	if res.RowsAffected() == 0 {
		return repo.ErrAlreadyExists
	}

	return nil
}

func (pg *postgres) GetPullReq(ctx context.Context, pullReqId string) (*domain.PullReq, error) {
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

	exists, err := prExists(ctx, tx, pullReqId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, repo.ErrNotFound
	}

	pr := &domain.PullReq{Id: pullReqId}

	err = tx.QueryRow(ctx, "SELECT name, author_id, status WHERE id = $1", pullReqId).Scan(&pr.Name, &pr.AuthorId, &pr.Status)
	if err != nil {
		return nil, fmt.Errorf("while querying and scanning pr: %w", err)
	}

	return pr, nil
}

func (pg *postgres) GetUsersToAssign(ctx context.Context, authorId string) ([]domain.User, error) {
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

	exists, err := userExists(ctx, tx, authorId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, repo.ErrNotFound
	}

	var teamName string
	if err := tx.QueryRow(ctx, "SELECT team_name FROM users WHERE id = $1", authorId).Scan(&teamName); err != nil {
		return nil, fmt.Errorf("getting team name of the author: %w", err)
	}

	users := []domain.User{}

	rows, err := tx.Query(
		ctx,
		`
			SELECT id, username, is_active FROM users
			WHERE team_name = $1 AND is_active = true
			LIMIT 2
		`,
		teamName,
	)
	if err != nil {
		return nil, fmt.Errorf("while querying team members: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.Id, &u.Username, &u.IsActive); err != nil {
			return nil, fmt.Errorf("while scanning user: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("while iterating on the rows: %w", err)
	}

	return users, nil
}
