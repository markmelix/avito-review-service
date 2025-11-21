package service

type Repo interface {
	TeamRepo
	UserRepo
	PullReqRepo
}
