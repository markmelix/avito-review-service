package pr

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"review/internal/domain"
	"review/internal/repo"
)

type PullReqRepo interface {
	CreatePullReq(ctx context.Context, pullReqId, name, authorId string) error
	GetPullReq(ctx context.Context, pullReqId string) (*domain.PullReq, error)
	GetUsersToAssign(ctx context.Context, authorId string) ([]domain.User, error)
	AssignPullReqReviewers(ctx context.Context, pullReqId string, reviewers []domain.User) error
	MarkPullReqMerged(ctx context.Context, pullReqId string) error
}

var (
	ErrPullReqAlreadyExists = errors.New("pr already exists")
	ErrAuthorNotFound       = errors.New("pr author not found")
	ErrPullReqNotFound      = errors.New("pr not found")
	ErrReviewerNotAssigned  = errors.New("reviewer is not assigned to this PR")
	ErrNoReassignCandidate  = errors.New("no active replacement candidate in team")
	ErrReassignOnMerged     = errors.New("cannot reassign on merged PR")
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

	users, err := uc.repo.GetUsersToAssign(ctx, authorId)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			return ErrAuthorNotFound
		default:
			return fmt.Errorf("while getting pr users to assign: %w", err)
		}
	}

	if err := uc.repo.AssignPullReqReviewers(ctx, pullReqId, users); err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			return ErrPullReqNotFound
		default:
			return fmt.Errorf("while assigning users to pr: %w", err)
		}
	}

	return nil
}

func (uc *PullReqService) MergePullReq(ctx context.Context, pullReqId string) error {
	if err := uc.repo.MarkPullReqMerged(ctx, pullReqId); err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			return ErrPullReqNotFound
		default:
			return fmt.Errorf("while merging pr: %w", err)
		}
	}
	return nil
}

func (uc *PullReqService) ReassignReviewer(ctx context.Context, pullReqId, oldReviewerId string) error {
	pr, err := uc.repo.GetPullReq(ctx, pullReqId)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			return ErrPullReqNotFound
		default:
			return fmt.Errorf("while getting pr: %w", err)
		}
	}

	if pr.Status == domain.PullReqMerged {
		return ErrReassignOnMerged
	}

	if !pr.HasReviewer(oldReviewerId) {
		return ErrReviewerNotAssigned
	}

	users, err := uc.repo.GetUsersToAssign(ctx, oldReviewerId)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			return ErrAuthorNotFound
		default:
			return fmt.Errorf("while getting pr users to assign: %w", err)
		}
	}

	if len(users) == 0 {
		return ErrNoReassignCandidate
	}

	user := users[rand.Intn(len(users))]

	if err := uc.repo.AssignPullReqReviewers(ctx, pullReqId, []domain.User{user}); err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			return ErrPullReqNotFound
		default:
			return fmt.Errorf("while assigning a user with id %v to pr: %w", user.Id, err)
		}
	}

	return nil
}
