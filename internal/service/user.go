package service

import (
	"context"
	"fmt"
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

type UserService struct {
	Repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo}
}

func (uc *UserService) SetIsActiveUser(ctx context.Context, id string, isActive bool) error {
	if err := uc.Repo.SetIsActiveUser(ctx, id, isActive); err != nil {
		return fmt.Errorf("failed toggling user: %w", err)
	}
	return nil
}

func (uc *UserService) GetReviewAssignments(ctx context.Context, id string) ([]PullReq, error) {
	prs, err := uc.Repo.GetReviewAssignments(ctx, id)
	if err != nil {
		return []PullReq{}, fmt.Errorf("failed getting user pr assignments: %w", err)
	}
	return prs, nil
}
