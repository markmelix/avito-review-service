package pr

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"review/internal/domain"
	"review/internal/repo"
)

const prAsigneeLimit int = 2

type PullReqRepo interface {
	CreatePullReq(ctx context.Context, pullReqId, name, authorId string) error
	GetPullReq(ctx context.Context, pullReqId string) (*domain.PullReq, error)
	GetUsersToAssign(ctx context.Context, authorId string, asigneeLimit int) ([]domain.User, error)
	AssignPullReqReviewers(ctx context.Context, pullReqId string, reviewers []domain.User, asigneeLimit int) error
	ReassignPullReqReviewer(ctx context.Context, pullReqId, oldReviewerId, newReviewerId string) error
	MarkPullReqMerged(ctx context.Context, pullReqId string) error
}

var (
	ErrPullReqAlreadyExists = errors.New("pr already exists")
	ErrAuthorNotFound       = errors.New("pr author not found")
	ErrPullReqNotFound      = errors.New("pr not found")
	ErrReviewerNotAssigned  = errors.New("reviewer is not assigned to this PR")
	ErrNoReassignCandidate  = errors.New("no active replacement candidate in team")
	ErrReassignOnMerged     = errors.New("cannot reassign on merged PR")
	ErrTooMuchAsignees      = errors.New("too much asignees")
)

type PullReqService struct {
	repo PullReqRepo
}

func NewPullReqService(repo PullReqRepo) *PullReqService {
	return &PullReqService{repo}
}

func (uc *PullReqService) CreatePullReq(ctx context.Context, pullReqId, name, authorId string) error {
	if err := uc.repo.CreatePullReq(ctx, pullReqId, name, authorId); err != nil {
		switch {
		case errors.Is(err, repo.ErrAlreadyExists):
			return ErrPullReqAlreadyExists
		case errors.Is(err, repo.ErrNotFound):
			return ErrAuthorNotFound
		default:
			return fmt.Errorf("while creating pull request: %w", err)
		}
	}

	users, err := uc.repo.GetUsersToAssign(ctx, authorId, prAsigneeLimit)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			return ErrAuthorNotFound
		default:
			return fmt.Errorf("while getting pr users to assign: %w", err)
		}
	}

	if err := uc.repo.AssignPullReqReviewers(ctx, pullReqId, users, prAsigneeLimit); err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			return ErrPullReqNotFound
		case errors.Is(err, repo.ErrAlreadyExists):
			return ErrTooMuchAsignees
		default:
			return fmt.Errorf("while assigning users to pr: %w", err)
		}
	}

	return nil
}

func (uc *PullReqService) MergePullReq(ctx context.Context, pullReqId string) (*domain.PullReq, error) {
	pr, err := uc.repo.GetPullReq(ctx, pullReqId)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			return nil, ErrPullReqNotFound
		default:
			return nil, fmt.Errorf("while getting pr: %w", err)
		}
	}

	if pr.Status == domain.PullReqMerged {
		return pr, nil
	}

	if err := uc.repo.MarkPullReqMerged(ctx, pullReqId); err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			return nil, ErrPullReqNotFound
		default:
			return nil, fmt.Errorf("while merging pr: %w", err)
		}
	}

	return pr, nil
}

// Reassigns reviewer with an active one from the team of the reviewer being
// replaced by and returns the PR and id of the user to be replaced with
func (uc *PullReqService) ReassignReviewer(ctx context.Context, pullReqId, oldReviewerId string) (*domain.PullReq, string, error) {
	pr, err := uc.repo.GetPullReq(ctx, pullReqId)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			return nil, "", ErrPullReqNotFound
		default:
			return nil, "", fmt.Errorf("while getting pr: %w", err)
		}
	}

	if pr.Status == domain.PullReqMerged {
		return nil, "", ErrReassignOnMerged
	}

	if !pr.HasReviewer(oldReviewerId) {
		return nil, "", ErrReviewerNotAssigned
	}

	users, err := uc.repo.GetUsersToAssign(ctx, oldReviewerId, prAsigneeLimit)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			return nil, "", ErrAuthorNotFound
		default:
			return nil, "", fmt.Errorf("while getting pr users to assign: %w", err)
		}
	}

	if len(users) == 0 {
		return nil, "", ErrNoReassignCandidate
	}

	user := users[rand.Intn(len(users))]

	if err := uc.repo.ReassignPullReqReviewer(ctx, pullReqId, oldReviewerId, user.Id); err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			return nil, "", ErrPullReqNotFound
		default:
			return nil, "", fmt.Errorf("while assigning a user with id %v to pr: %w", user.Id, err)
		}
	}

	return pr, user.Id, nil
}

func (uc *PullReqService) GetPullReq(ctx context.Context, pullReqId string) (*domain.PullReq, error) {
	pr, err := uc.repo.GetPullReq(ctx, pullReqId)

	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			return nil, ErrPullReqNotFound
		default:
			return nil, fmt.Errorf("while getting pr: %w", err)
		}
	}

	return pr, nil
}
