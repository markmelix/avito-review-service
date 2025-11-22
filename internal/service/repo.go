package service

import (
	"review/internal/service/pr"
	"review/internal/service/team"
	"review/internal/service/user"
)

type Repo interface {
	team.TeamRepo
	user.UserRepo
	pr.PullReqRepo
}
