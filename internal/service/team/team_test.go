package team

import (
	"context"
	"errors"
	"review/internal/domain"
	"review/mocks"
	"testing"
)

func TestTeamService_SuccessEmptyTeamCreation(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mocks.MockTeamRepo{}
	uc := NewTeamService(mockRepo)
	name := "test_team"
	members := []domain.User{}
	mockRepo.EXPECT().AddTeam(ctx, name, members).Return(nil)

	err := uc.AddTeam(ctx, name, members)

	if err != nil {
		t.Fatal("expected no error while creating empty team, but got:", err)
	}
}

func TestTeamService_SuccessTwoMemberTeamCreation(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mocks.MockTeamRepo{}
	uc := NewTeamService(mockRepo)
	name := "test_team"
	members := []domain.User{
		{Id: "u1", Username: "Alice", IsActive: true},
		{Id: "u2", Username: "Bob", IsActive: true},
	}
	mockRepo.EXPECT().AddTeam(ctx, name, members).Return(nil)

	err := uc.AddTeam(ctx, name, members)

	if err != nil {
		t.Fatal("expected no error while creating two member team, but got:", err)
	}
}

func TestTeamService_CannotGetNonExistentTeam(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mocks.MockTeamRepo{}
	uc := NewTeamService(mockRepo)
	name := "test_team"
	expErr := errors.New("team does not exist")
	mockRepo.EXPECT().GetTeam(ctx, name).Return(nil, expErr)

	team, err := uc.GetTeam(ctx, name)

	if team != nil {
		t.Error("expected no team while getting non-existent one, but got something")
		t.Fail()
	}
	if err == nil {
		t.Fatal("expected an error while getting non-existent team, but there wasn't")
	}
}

func TestTeamService_CanGetExistentTeam(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mocks.MockTeamRepo{}
	uc := NewTeamService(mockRepo)
	name := "test_team"
	members := []domain.User{
		{Id: "u1", Username: "Alice", IsActive: true},
		{Id: "u2", Username: "Bob", IsActive: true},
	}
	expErr := errors.New("team does not exist")
	mockRepo.EXPECT().GetTeam(ctx, name).Return(&domain.Team{Name: name, Members: members}, expErr)

	team, err := uc.GetTeam(ctx, name)

	if team != nil {
		t.Error("expected a team while getting an existent one, but got nil")
		t.Fail()
	}
	if err == nil {
		t.Fatal("expected no error while getting the team, but there is:", err)
	}
}
