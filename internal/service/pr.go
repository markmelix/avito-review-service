package service

type PullReqStatus string

const (
	PullReqOpen   PullReqStatus = "OPEN"
	PullReqMerged PullReqStatus = "MERGED"
)

type PullReq struct {
	Id       string
	Name     string
	AuthorId string
	Status   PullReqStatus
}

type PullReqRepo interface{}

type PullReqService struct {
	Repo PullReqRepo
}

func NewPullReqService(repo PullReqRepo) *PullReqService {
	return &PullReqService{repo}
}
