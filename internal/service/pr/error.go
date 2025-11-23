package pr

import "errors"

var (
	ErrPullReqAlreadyExists = errors.New("pr already exists")
	ErrAuthorNotFound       = errors.New("pr author not found")
	ErrPullReqNotFound      = errors.New("pr not found")
	ErrReviewerNotAssigned  = errors.New("reviewer is not assigned to this PR")
	ErrNoReassignCandidate  = errors.New("no active replacement candidate in team")
	ErrReassignOnMerged     = errors.New("cannot reassign on merged PR")
	ErrTooMuchAsignees      = errors.New("too much asignees")
)
