package domain

import "slices"

type PullReqStatus string

const (
	PullReqOpen   PullReqStatus = "OPEN"
	PullReqMerged PullReqStatus = "MERGED"
)

type PullReq struct {
	Id        string
	Name      string
	AuthorId  string
	Status    PullReqStatus
	Reviewers []User
}

func (pr PullReq) HasReviewer(id string) bool {
	return slices.ContainsFunc(pr.Reviewers, func(u User) bool {
		return u.Id == id
	})
}
