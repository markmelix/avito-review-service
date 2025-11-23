package pr

import "review/internal/domain"

type PullReqCreateReqDto struct {
	PullReqId   string `json:"pull_request_id"`
	PullReqName string `json:"pull_request_name"`
	AuthorId    string `json:"author_id"`
}

type PullReqCreateRespBody struct {
	PullReqId   string               `json:"pull_request_id"`
	PullReqName string               `json:"pull_request_name"`
	AuthorId    string               `json:"author_id"`
	Status      domain.PullReqStatus `json:"status"`
	ReviewerIds []string             `json:"assigned_reviewers"`
}

type PullReqCreateRespDto struct {
	PullReq PullReqCreateRespBody `json:"pr"`
}

func usersToIds(users []domain.User) []string {
	s := make([]string, len(users))
	for i, u := range users {
		s[i] = u.Id
	}
	return s
}

func CreateRespDtoFromPullReq(pr *domain.PullReq) *PullReqCreateRespDto {
	return &PullReqCreateRespDto{
		PullReq: PullReqCreateRespBody{
			PullReqId:   pr.Id,
			PullReqName: pr.Name,
			AuthorId:    pr.AuthorId,
			Status:      pr.Status,
			ReviewerIds: usersToIds(pr.Reviewers),
		},
	}
}
