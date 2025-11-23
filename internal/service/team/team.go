package team

import (
	"context"
	"errors"
	"fmt"
	"review/internal/domain"
	"review/internal/repo"
)

type TeamRepo interface {
	AddTeam(ctx context.Context, name string, members []domain.User) error
	GetTeam(ctx context.Context, name string) (*domain.Team, error)
}

type TeamService struct {
	repo TeamRepo
}

func NewTeamService(repo TeamRepo) *TeamService {
	return &TeamService{repo}
}

func (uc *TeamService) AddTeam(ctx context.Context, name string, members []domain.User) error {
	if err := uc.repo.AddTeam(ctx, name, members); err != nil {
		if errors.Is(err, repo.ErrAlreadyExists) {
			return ErrTeamAlreadyExists
		}
		return fmt.Errorf("failed creating team with members: %w", err)
	}
	return nil
}

func (uc *TeamService) GetTeam(ctx context.Context, name string) (*domain.Team, error) {
	team, err := uc.repo.GetTeam(ctx, name)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrTeamNotFound
		}
		return nil, fmt.Errorf("failed creating team with members: %w", err)
	}
	return team, nil
}
