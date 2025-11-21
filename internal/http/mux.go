package handler

import (
	"net/http"
	"review/internal/http/handler"
	"review/internal/service"
	"strings"
)

func joinMux(root *http.ServeMux, inner http.Handler, prefix string) {
	prefix = strings.TrimRight(prefix, "/")
	root.Handle(prefix+"/", http.StripPrefix(prefix, inner))
}

func NewMux(repo service.Repo) http.Handler {
	teamMux := handler.NewTeamMux(handler.NewTeamHandler(service.NewTeamService(repo)))
	userMux := handler.NewUserMux(handler.NewUserHandler(service.NewUserService(repo)))
	prMux := handler.NewPullReqMux(handler.NewPullReqHandler(service.NewPullReqService(repo)))

	mux := http.NewServeMux()
	joinMux(mux, teamMux, "/team")
	joinMux(mux, userMux, "/users")
	joinMux(mux, prMux, "/pullRequest")
	return corsMiddleware(mux)
}
