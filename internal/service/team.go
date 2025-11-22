package service

import "context"

type Team struct {
	Name    string
	Members []User
}

type TeamRepo interface {
	AddTeam(ctx context.Context, name string, members []User) error
}

type TeamService struct {
	Repo TeamRepo
}

func NewTeamService(repo TeamRepo) *TeamService {
	return &TeamService{repo}
}
