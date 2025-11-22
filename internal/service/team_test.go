package service_test

import (
	"context"
	"errors"
	"review/internal/service"
	"review/mocks"
	"testing"
)

func TestTeamService_SuccessEmptyTeamCreation(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mocks.MockTeamRepo{}
	uc := service.NewTeamService(mockRepo)
	name := "test_team"
	members := []service.User{}
	mockRepo.EXPECT().AddTeam(ctx, name, members).Return(nil)

	err := uc.AddTeam(ctx, name, members)

	if err != nil {
		t.Fatal("expected no error while creating empty team, but got:", err)
	}
}

func TestTeamService_SuccessTwoMemberTeamCreation(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mocks.MockTeamRepo{}
	uc := service.NewTeamService(mockRepo)
	name := "test_team"
	members := []service.User{{"u1", "Alice", true}, {"u2", "Bob", true}}
	mockRepo.EXPECT().AddTeam(ctx, name, members).Return(nil)

	err := uc.AddTeam(ctx, name, members)

	if err != nil {
		t.Fatal("expected no error while creating two member team, but got:", err)
	}
}

func TestTeamService_CannotGetNonExistentTeam(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mocks.MockTeamRepo{}
	uc := service.NewTeamService(mockRepo)
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
	uc := service.NewTeamService(mockRepo)
	name := "test_team"
	members := []service.User{{"u1", "Alice", true}, {"u2", "Bob", true}}
	expErr := errors.New("team does not exist")
	mockRepo.EXPECT().GetTeam(ctx, name).Return(&service.Team{name, members}, expErr)

	team, err := uc.GetTeam(ctx, name)

	if team != nil {
		t.Error("expected a team while getting an existent one, but got nil")
		t.Fail()
	}
	if err == nil {
		t.Fatal("expected no error while getting the team, but there is:", err)
	}
}
