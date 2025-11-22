package service

type User struct {
	Id       string
	Username string
	IsActive bool
}

type UserRepo interface{}

type UserService struct {
	Repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo}
}
