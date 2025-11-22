package user

import (
	"context"
	"review/internal/domain"
	"review/internal/repo"
	"review/mocks"
	"testing"
)

func TestUserService_SetIsActiveAlertsNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mocks.MockUserRepo{}
	uc := NewUserService(mockRepo)
	id := "u1"
	isActive := false
	mockRepo.EXPECT().SetIsActiveUser(ctx, id, isActive).Return(repo.ErrNotFound)

	err := uc.SetIsActiveUser(ctx, id, isActive)

	if err != ErrUserNotFound {
		t.Fatalf("expected: \"%v\" error, got: %v", ErrUserNotFound, err)
	}
}

func TestUserService_SetIsActiveExistingUserSuccessfully(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mocks.MockUserRepo{}
	uc := NewUserService(mockRepo)
	id := "u1"
	isActive := false
	mockRepo.EXPECT().SetIsActiveUser(ctx, id, isActive).Return(nil)

	err := uc.SetIsActiveUser(ctx, id, isActive)

	if err != nil {
		t.Fatal("expected no error, got:", err)
	}
}

func TestUserService_GetReviewAssignmentsAlertsNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mocks.MockUserRepo{}
	uc := NewUserService(mockRepo)
	id := "u1"
	mockRepo.EXPECT().GetReviewAssignments(ctx, id).Return(nil, repo.ErrNotFound)

	prs, err := uc.GetReviewAssignments(ctx, id)

	if prs != nil {
		t.Error("first value returned is expected to be nil, got:", prs)
		t.Fail()
	}
	if err != ErrUserNotFound {
		t.Fatalf("expected: \"%v\" error, got: %v", ErrUserNotFound, err)
	}
}

func TestUserService_GetReviewAssignmentsSuccess(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mocks.MockUserRepo{}
	uc := NewUserService(mockRepo)
	id := "u1"
	expPrs := []domain.PullReq{{
		Id:        "pr-1001",
		Name:      "Add search",
		AuthorId:  "u1",
		Status:    domain.PullReqOpen,
		Reviewers: nil,
	}}
	mockRepo.EXPECT().GetReviewAssignments(ctx, id).Return(expPrs, nil)

	prs, err := uc.GetReviewAssignments(ctx, id)

	if err != nil {
		t.Error("expected no error, got:", err)
		t.Fail()
	}
	if prs == nil {
		t.Fatal("pr slice is not expected to be nil")
	}
	if len(prs) != len(expPrs) {
		t.Fatal("pr slice expected to be of length one would expect, but it is not")
	}
	if prs[0].Id != expPrs[0].Id {
		t.Fatalf("expected: %v, got: %v", expPrs[0].Id, prs[0].Id)
	}
}
