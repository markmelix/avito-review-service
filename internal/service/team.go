package service

import (
	"context"
	"fmt"
)

type Team struct {
	Name    string
	Members []User
}

type TeamRepo interface {
	AddTeam(ctx context.Context, name string, members []User) error
	GetTeam(ctx context.Context, name string) (*Team, error)
}

type TeamService struct {
	Repo TeamRepo
}

func NewTeamService(repo TeamRepo) *TeamService {
	return &TeamService{repo}
}

func (uc *TeamService) AddTeam(ctx context.Context, name string, members []User) error {
	if err := uc.Repo.AddTeam(ctx, name, members); err != nil {
		return fmt.Errorf("failed creating team with members: %w", err)
	}
	return nil
}

func (uc *TeamService) GetTeam(ctx context.Context, name string) (*Team, error) {
	team, err := uc.Repo.GetTeam(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed creating team with members: %w", err)
	}
	return team, nil
}
