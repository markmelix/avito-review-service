package user

import (
	"context"
	"errors"
	"fmt"
	"review/internal/domain"
	"review/internal/repo"
)

type UserRepo interface {
	SetIsActiveUser(ctx context.Context, id string, isActive bool) error
	GetReviewAssignments(ctx context.Context, id string) ([]domain.PullReq, error)
}

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo}
}

func (uc *UserService) SetIsActiveUser(ctx context.Context, id string, isActive bool) error {
	if err := uc.repo.SetIsActiveUser(ctx, id, isActive); err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("toggling user: %w", err)
	}
	return nil
}

func (uc *UserService) GetReviewAssignments(ctx context.Context, id string) ([]domain.PullReq, error) {
	prs, err := uc.repo.GetReviewAssignments(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("getting user pr assignments: %w", err)
	}
	return prs, nil
}
