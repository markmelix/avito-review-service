package postgres

import (
	"context"
	"errors"
	"fmt"
	"review/internal/domain"
	"review/internal/repo"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return repo.ErrAlreadyExists
		} else if pgErr != nil && pgErr.Code == "23503" {
			return repo.ErrNotFound
		}
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

	err = tx.QueryRow(ctx, "SELECT name, author_id, status, merged_at FROM pull_requests WHERE id = $1", pullReqId).Scan(&pr.Name, &pr.AuthorId, &pr.Status, &pr.MergedAt)
	if err != nil {
		return nil, fmt.Errorf("while querying and scanning pr: %w", err)
	}

	rows, err := tx.Query(
		ctx,
		`
			SELECT u.id, u.username, u.is_active, t.name AS team_name
			FROM users u
			JOIN teams t ON t.id = u.team_id
			WHERE u.id IN (SELECT user_id FROM assignments WHERE pr_id = $1)
		`,
		pullReqId,
	)
	if err != nil {
		return nil, fmt.Errorf("while querying team members: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.Id, &u.Username, &u.IsActive, &u.TeamName); err != nil {
			return nil, fmt.Errorf("while scanning user: %w", err)
		}
		pr.Reviewers = append(pr.Reviewers, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("while iterating on the rows: %w", err)
	}

	return pr, nil
}

func (pg *postgres) GetUsersToAssign(ctx context.Context, authorId string, asigneeLimit int, pullReqId *string) ([]domain.User, error) {
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
	if err := tx.QueryRow(ctx, "SELECT t.name FROM users u JOIN teams t ON t.id = u.team_id WHERE u.id = $1", authorId).Scan(&teamName); err != nil {
		return nil, fmt.Errorf("getting team name of the author: %w", err)
	}

	users := []domain.User{}

	prId := ""
	if pullReqId != nil {
		prId = *pullReqId
	}

	rows, err := tx.Query(
		ctx,
		`
			SELECT u.id, u.username, u.is_active
			FROM users u
			JOIN teams t ON t.id = u.team_id
			WHERE t.name = $1 AND
				u.is_active = true AND
				u.id != $3 AND
				($4 = '' OR u.id NOT IN
					(SELECT user_id FROM assignments WHERE pr_id = $4))
			LIMIT $2
		`,
		teamName,
		asigneeLimit,
		authorId,
		prId,
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

func (pg *postgres) AssignPullReqReviewers(ctx context.Context, pullReqId string, reviewers []domain.User, asigneeLimit int) error {
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

	exists, err := prExists(ctx, tx, pullReqId)
	if err != nil {
		return err
	}
	if !exists {
		return repo.ErrNotFound
	}

	if len(reviewers) > asigneeLimit {
		return repo.ErrAlreadyExists
	}

	for _, r := range reviewers {
		tx.Exec(ctx, "INSERT INTO assignments (user_id, pr_id) VALUES ($1, $2)", r.Id, pullReqId)
	}

	var asigneeAmount int
	tx.QueryRow(ctx, "SELECT COUNT(*) FROM assignments WHERE pr_id = $1", pullReqId).Scan(&asigneeAmount)
	if asigneeAmount > asigneeLimit {
		return repo.ErrAlreadyExists
	}

	return nil
}

func (pg *postgres) ReassignPullReqReviewer(ctx context.Context, pullReqId, oldReviewerId, newReviewerId string) error {
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

	exists, err := prExists(ctx, tx, pullReqId)
	if err != nil {
		return err
	}
	if !exists {
		return repo.ErrNotFound
	}

	_, err = tx.Exec(ctx, "UPDATE assignments SET user_id = $1 WHERE user_id = $2 AND pr_id = $3", newReviewerId, oldReviewerId, pullReqId)
	if err != nil {
		return fmt.Errorf("updating assignments with new user_id: %w", err)
	}

	return nil
}

func (pg *postgres) MarkPullReqMerged(ctx context.Context, pullReqId string) error {
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

	exists, err := prExists(ctx, tx, pullReqId)
	if err != nil {
		return err
	}
	if !exists {
		return repo.ErrNotFound
	}

	_, err = tx.Exec(ctx, "UPDATE pull_requests SET status = $1, merged_at = NOW() WHERE id = $2", domain.PullReqMerged, pullReqId)
	if err != nil {
		return fmt.Errorf("updating pr making it merged: %w", err)
	}

	return nil
}
