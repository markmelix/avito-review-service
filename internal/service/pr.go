package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"review/internal/repo"
)

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

type PullReqRepo interface {
	CreatePullReq(ctx context.Context, pullReqId, name, authorId string) error
	GetUsersToAssign(ctx context.Context, authorId string) ([]User, error)
	AssignPullReqReviewers(ctx context.Context, pullReqId string, reviewers []User) error
	MarkPullReqMerged(ctx context.Context, pullReqId string) error
}

var (
	ErrPullReqAlreadyExists = errors.New("pr already exists")
	ErrAuthorNotFound       = errors.New("pr author not found")
	ErrPullReqNotFound      = errors.New("pr not found")
)

type PullReqService struct {
	Repo PullReqRepo
}

func NewPullReqService(repo PullReqRepo) *PullReqService {
	return &PullReqService{repo}
}

func (uc *PullReqService) CreatePullReq(ctx context.Context, pullReqId, name, authorId string) error {
	if err := uc.Repo.CreatePullReq(ctx, pullReqId, name, authorId); err != nil {
		switch {
		case errors.Is(err, repo.ErrAlreadyExists):
			return ErrPullReqAlreadyExists
		case errors.Is(err, repo.ErrNotFound):
			return ErrAuthorNotFound
		default:
			return fmt.Errorf("while creating pull request: %w", err)
		}
	}

	users, err := uc.Repo.GetUsersToAssign(ctx, authorId)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			return ErrAuthorNotFound
		default:
			return fmt.Errorf("while getting pr users to assign: %w", err)
		}
	}

	if err := uc.Repo.AssignPullReqReviewers(ctx, pullReqId, users); err != nil {
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
	if err := uc.Repo.MarkPullReqMerged(ctx, pullReqId); err != nil {
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
	users, err := uc.Repo.GetUsersToAssign(ctx, oldReviewerId)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			return ErrAuthorNotFound
		default:
			return fmt.Errorf("while getting pr users to assign: %w", err)
		}
	}
	if len(users) == 0 {
		return nil
	}

	user := users[rand.Intn(len(users))]
	if err := uc.Repo.AssignPullReqReviewers(ctx, pullReqId, []User{user}); err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			return ErrPullReqNotFound
		default:
			return fmt.Errorf("while assigning a user with id %v to pr: %w", user.Id, err)
		}
	}

	return nil
}
