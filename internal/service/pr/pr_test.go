package pr

import (
	"context"
	"errors"
	"review/internal/domain"
	"review/mocks"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestPullReqService_CreatePullReq_WontAssignMoreThanTwoReviewers(t *testing.T) {
	var (
		ctx      = context.Background()
		mockRepo = &mocks.MockPullReqRepo{}
		uc       = NewPullReqService(mockRepo)
		prRevLim = 2
		prId     = "pr-1001"
		prName   = "Add search"
		autId    = "u1"
	)

	mockRepo.
		On("CreatePullReq", mock.Anything, prId, prName, autId).
		Return(nil)
	mockRepo.
		On("GetUsersToAssign", mock.Anything, autId, prRevLim, mock.Anything).
		Return([]domain.User{{}, {}, {}}, nil)

	err := uc.CreatePullReq(ctx, prId, prName, autId)

	if !errors.Is(err, ErrTooMuchAsignees) {
		t.Fatalf("expected %v error, got: %v", ErrTooMuchAsignees, err)
	}

	mockRepo.AssertExpectations(t)
}

func TestPullReqService_CreatePullReq_WillAssignReviewers(t *testing.T) {
	var (
		ctx      = context.Background()
		mockRepo = &mocks.MockPullReqRepo{}
		uc       = NewPullReqService(mockRepo)
		users    = []domain.User{{}, {}}
		prRevLim = 2
		prId     = "pr-1001"
		prName   = "Add search"
		autId    = "u1"
	)

	mockRepo.
		On("CreatePullReq", mock.Anything, prId, prName, autId).
		Return(nil)
	mockRepo.
		On("GetUsersToAssign", mock.Anything, autId, prRevLim, mock.Anything).
		Return(users, nil)
	mockRepo.
		On("AssignPullReqReviewers", mock.Anything, prId, users, prRevLim).
		Return(nil)

	err := uc.CreatePullReq(ctx, prId, prName, autId)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	mockRepo.AssertExpectations(t)
}

func TestPullReqService_MergePullReq_WillCallMarkMergedOnOpenedPR(t *testing.T) {
	var (
		ctx      = context.Background()
		mockRepo = &mocks.MockPullReqRepo{}
		uc       = NewPullReqService(mockRepo)
		pr       = &domain.PullReq{Status: domain.PullReqOpen}
		prId     = "pr-1001"
	)

	mockRepo.
		On("GetPullReq", mock.Anything, prId).
		Return(pr, nil).
		Twice()
	mockRepo.
		On("MarkPullReqMerged", mock.Anything, prId).
		Return(nil).
		Once()

	_, err := uc.MergePullReq(ctx, prId)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	mockRepo.AssertExpectations(t)
}

func TestPullReqService_MergePullReq_WontCallMarkMergedOnAlreadyMergedPR(t *testing.T) {
	var (
		ctx      = context.Background()
		mockRepo = &mocks.MockPullReqRepo{}
		uc       = NewPullReqService(mockRepo)
		pr       = &domain.PullReq{Status: domain.PullReqMerged}
		prId     = "pr-1001"
	)

	mockRepo.
		On("GetPullReq", mock.Anything, prId).
		Return(pr, nil).
		Once()

	_, err := uc.MergePullReq(ctx, prId)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	mockRepo.AssertExpectations(t)
}

func TestPullReqService_ReassignReviewer_WontReassignOnMergedPR(t *testing.T) {
	var (
		ctx      = context.Background()
		mockRepo = &mocks.MockPullReqRepo{}
		uc       = NewPullReqService(mockRepo)
		pr       = &domain.PullReq{Status: domain.PullReqMerged}
		prId     = "pr-1001"
		revId    = "u1"
	)

	mockRepo.
		On("GetPullReq", mock.Anything, prId).
		Return(pr, nil).
		Once()

	_, _, err := uc.ReassignReviewer(ctx, prId, revId)

	if !errors.Is(err, ErrReassignOnMerged) {
		t.Fatalf("expected %v error, got: %v", ErrReassignOnMerged, err)
	}

	mockRepo.AssertExpectations(t)
}

func TestPullReqService_ReassignReviewer_WontReassignUnassignedReviewer(t *testing.T) {
	var (
		ctx      = context.Background()
		mockRepo = &mocks.MockPullReqRepo{}
		uc       = NewPullReqService(mockRepo)
		pr       = &domain.PullReq{Status: domain.PullReqOpen, Reviewers: []domain.User{}}
		prId     = "pr-1001"
		revId    = "u1"
	)

	mockRepo.
		On("GetPullReq", mock.Anything, prId).
		Return(pr, nil).
		Once()

	_, _, err := uc.ReassignReviewer(ctx, prId, revId)

	if !errors.Is(err, ErrReviewerNotAssigned) {
		t.Fatalf("expected %v error, got: %v", ErrReviewerNotAssigned, err)
	}

	mockRepo.AssertExpectations(t)
}

func TestPullReqService_ReassignReviewer_WontReassignZeroCandidates(t *testing.T) {
	var (
		ctx      = context.Background()
		mockRepo = &mocks.MockPullReqRepo{}
		uc       = NewPullReqService(mockRepo)
		revId    = "u1"
		prId     = "pr-1001"
		users    = []domain.User{}
		pr       = &domain.PullReq{
			Status:    domain.PullReqOpen,
			Reviewers: []domain.User{{Id: revId}},
		}
	)

	mockRepo.
		On("GetPullReq", mock.Anything, prId).
		Return(pr, nil).
		Once()
	mockRepo.
		On("GetUsersToAssign", mock.Anything, revId, prAsigneeLimit, &prId).
		Return(users, nil).
		Once()

	_, _, err := uc.ReassignReviewer(ctx, prId, revId)

	if !errors.Is(err, ErrNoReassignCandidate) {
		t.Fatalf("expected %v error, got: %v", ErrNoReassignCandidate, err)
	}

	mockRepo.AssertExpectations(t)
}

func TestPullReqService_ReassignReviewer_ReassignsSuccessfully(t *testing.T) {
	var (
		ctx      = context.Background()
		mockRepo = &mocks.MockPullReqRepo{}
		uc       = NewPullReqService(mockRepo)
		revId    = "u1"
		newId    = "u2"
		prId     = "pr-1001"
		users    = []domain.User{{Id: newId}}
		pr       = &domain.PullReq{
			Status:    domain.PullReqOpen,
			Reviewers: []domain.User{{Id: revId}},
		}
	)

	mockRepo.
		On("GetPullReq", mock.Anything, prId).
		Return(pr, nil).
		Twice()
	mockRepo.
		On("GetUsersToAssign", mock.Anything, revId, prAsigneeLimit, &prId).
		Return(users, nil).
		Once()
	mockRepo.
		On("ReassignPullReqReviewer", mock.Anything, prId, revId, newId).
		Return(nil).
		Once()

	_, _, err := uc.ReassignReviewer(ctx, prId, revId)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	mockRepo.AssertExpectations(t)
}
