package service

type PullReqRepo interface{}

type PullReqService struct {
	Repo PullReqRepo
}

func NewPullReqService(repo PullReqRepo) *PullReqService {
	return &PullReqService{repo}
}
