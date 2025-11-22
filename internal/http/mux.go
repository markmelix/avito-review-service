package handler

import (
	"net/http"
	"review/internal/http/handler"
	"review/internal/service"
	"review/internal/service/pr"
	"review/internal/service/team"
	"review/internal/service/user"
	"strings"
)

func joinMux(root *http.ServeMux, inner http.Handler, prefix string) {
	prefix = strings.TrimRight(prefix, "/")
	root.Handle(prefix+"/", http.StripPrefix(prefix, inner))
}

func NewMux(repo service.Repo) http.Handler {
	teamMux := handler.NewTeamMux(handler.NewTeamHandler(team.NewTeamService(repo)))
	userMux := handler.NewUserMux(handler.NewUserHandler(user.NewUserService(repo)))
	prMux := handler.NewPullReqMux(handler.NewPullReqHandler(pr.NewPullReqService(repo)))

	mux := http.NewServeMux()
	joinMux(mux, teamMux, "/team")
	joinMux(mux, userMux, "/users")
	joinMux(mux, prMux, "/pullRequest")
	return corsMiddleware(mux)
}
