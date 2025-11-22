package service

import (
	"context"
	"errors"
	"fmt"
	"review/internal/repo"
)

type User struct {
	Id       string
	Username string
	IsActive bool
}

type UserRepo interface {
	SetIsActiveUser(ctx context.Context, id string, isActive bool) error
	GetReviewAssignments(ctx context.Context, id string) ([]PullReq, error)
}

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
	Repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo}
}

func (uc *UserService) SetIsActiveUser(ctx context.Context, id string, isActive bool) error {
	if err := uc.Repo.SetIsActiveUser(ctx, id, isActive); err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("toggling user: %w", err)
	}
	return nil
}

func (uc *UserService) GetReviewAssignments(ctx context.Context, id string) ([]PullReq, error) {
	prs, err := uc.Repo.GetReviewAssignments(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("getting user pr assignments: %w", err)
	}
	return prs, nil
}
