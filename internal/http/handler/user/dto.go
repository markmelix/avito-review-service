package user

import "review/internal/domain"

type SetIsActiveDto struct {
	UserId   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type UserDto struct {
	User domain.User `json:"user"`
}

type GetReviewDto struct {
	UserId   string           `json:"user_id"`
	PullReqs []domain.PullReq `json:"pull_requests"`
}
