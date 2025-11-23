package pr

import (
	"review/internal/domain"
	"time"
)

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

type PullReqMergeReqDto struct {
	PullReqId string `json:"pull_request_id"`
}

type PullReqMergeRespBody struct {
	PullReqId   string               `json:"pull_request_id"`
	PullReqName string               `json:"pull_request_name"`
	AuthorId    string               `json:"author_id"`
	Status      domain.PullReqStatus `json:"status"`
	ReviewerIds []string             `json:"assigned_reviewers"`
	MergedAt    string               `json:"mergedAt"`
}

type PullReqMergeRespDto struct {
	PullReq PullReqMergeRespBody `json:"pr"`
}

func MergeRespDtoFromPullReq(pr *domain.PullReq) *PullReqMergeRespDto {
	return &PullReqMergeRespDto{
		PullReq: PullReqMergeRespBody{
			PullReqId:   pr.Id,
			PullReqName: pr.Name,
			AuthorId:    pr.AuthorId,
			Status:      pr.Status,
			ReviewerIds: usersToIds(pr.Reviewers),
			MergedAt:    pr.MergedAt.Format(time.RFC3339),
		},
	}
}

type PullReqReassignReqDto struct {
	PullReqId     string `json:"pull_request_id"`
	OldReviewerId string `json:"old_reviewer_id"`
}

type PullReqReassignRespDto struct {
	PullReq      PullReqCreateRespBody `json:"pr"`
	ReplacedById string                `json:"replaced_by"`
}

func ReassignRespFromPullReq(pr *domain.PullReq, replacedById string) *PullReqReassignRespDto {
	return &PullReqReassignRespDto{
		PullReq:      CreateRespDtoFromPullReq(pr).PullReq,
		ReplacedById: replacedById,
	}
}
