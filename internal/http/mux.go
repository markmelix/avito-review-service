package handler

import (
	"net/http"
	"review/internal/http/handler/pr"
	"review/internal/http/handler/team"
	"review/internal/http/handler/user"
	"review/internal/service"
	pr_uc "review/internal/service/pr"
	team_uc "review/internal/service/team"
	user_uc "review/internal/service/user"
	"strings"
)

func joinMux(root *http.ServeMux, inner http.Handler, prefix string) {
	prefix = strings.TrimRight(prefix, "/")
	root.Handle(prefix+"/", http.StripPrefix(prefix, inner))
}

func NewMux(repo service.Repo) http.Handler {
	teamMux := team.NewTeamMux(team.NewTeamHandler(team_uc.NewTeamService(repo)))
	userMux := user.NewUserMux(user.NewUserHandler(user_uc.NewUserService(repo)))
	prMux := pr.NewPullReqMux(pr.NewPullReqHandler(pr_uc.NewPullReqService(repo)))

	mux := http.NewServeMux()
	joinMux(mux, teamMux, "/team")
	joinMux(mux, userMux, "/users")
	joinMux(mux, prMux, "/pullRequest")
	return corsMiddleware(mux)
}
