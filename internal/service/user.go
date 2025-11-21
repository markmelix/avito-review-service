package service

type UserRepo interface{}

type UserService struct {
	Repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo}
}
