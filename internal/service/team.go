package service

import (
	"context"
	"errors"
	"fmt"
	"review/internal/repo"
)

type Team struct {
	Name    string
	Members []User
}

type TeamRepo interface {
	AddTeam(ctx context.Context, name string, members []User) error
	GetTeam(ctx context.Context, name string) (*Team, error)
}

var (
	ErrTeamAlreadyExists = errors.New("team already exists")
	ErrTeamNotFound      = errors.New("team not found")
)

type TeamService struct {
	repo TeamRepo
}

func NewTeamService(repo TeamRepo) *TeamService {
	return &TeamService{repo}
}

func (uc *TeamService) AddTeam(ctx context.Context, name string, members []User) error {
	if err := uc.repo.AddTeam(ctx, name, members); err != nil {
		if errors.Is(err, repo.ErrAlreadyExists) {
			return ErrTeamAlreadyExists
		}
		return fmt.Errorf("failed creating team with members: %w", err)
	}
	return nil
}

func (uc *TeamService) GetTeam(ctx context.Context, name string) (*Team, error) {
	team, err := uc.repo.GetTeam(ctx, name)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrTeamNotFound
		}
		return nil, fmt.Errorf("failed creating team with members: %w", err)
	}
	return team, nil
}
