package service

type TeamRepo interface{}

type TeamService struct {
	Repo TeamRepo
}

func NewTeamService(repo TeamRepo) *TeamService {
	return &TeamService{repo}
}
