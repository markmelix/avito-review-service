package domain

import (
	"slices"
	"time"
)

type PullReqStatus string

const (
	PullReqOpen   PullReqStatus = "OPEN"
	PullReqMerged PullReqStatus = "MERGED"
)

type PullReq struct {
	Id        string        `json:"pull_request_id"`
	Name      string        `json:"pull_request_name"`
	AuthorId  string        `json:"author_id"`
	Status    PullReqStatus `json:"status"`
	Reviewers []User        `json:"-"`
	MergedAt  *time.Time    `json:"-"`
}

func (pr PullReq) HasReviewer(id string) bool {
	return slices.ContainsFunc(pr.Reviewers, func(u User) bool {
		return u.Id == id
	})
}
