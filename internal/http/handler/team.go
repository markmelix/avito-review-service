package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"review/internal/domain"
	"review/internal/http/helper"
	"review/internal/service/team"
)

type TeamHandler struct {
	usecase *team.TeamService
}

func NewTeamHandler(usecase *team.TeamService) *TeamHandler {
	return &TeamHandler{usecase}
}

type TeamAddDto struct {
	Team domain.Team `json:"team"`
}

func (h *TeamHandler) add(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		helper.WriteMethodNotAllowedError(w)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		helper.WriteUndefinedError(w, fmt.Errorf("could not read request body: %w", err))
		return
	}
	defer r.Body.Close()

	var t domain.Team
	if err := json.Unmarshal(body, &t); err != nil {
		helper.WriteUndefinedError(w, fmt.Errorf("could not unmarshall request body: %w", err))
		return
	}

	ctx := context.Background()
	if err := h.usecase.AddTeam(ctx, t.Name, t.Members); err != nil {
		switch {
		case errors.Is(err, team.ErrTeamAlreadyExists):
			helper.WriteError(w, "TEAM_EXISTS", "team_name already exists", http.StatusBadRequest)
		default:
			helper.WriteUndefinedError(w, err)
		}
		return
	}

	resp := &TeamAddDto{t}
	respJson, err := json.Marshal(resp)
	if err != nil {
		helper.WriteUndefinedError(w, fmt.Errorf("could not convert response to json: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}

func (h *TeamHandler) get(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		helper.WriteUndefinedError(w, errors.New("team_name query param was not specified"))
		return
	}

	ctx := context.Background()
	t, err := h.usecase.GetTeam(ctx, teamName)
	if err != nil {
		switch {
		case errors.Is(err, team.ErrTeamNotFound):
			helper.WriteNotFoundError(w)
		default:
			helper.WriteUndefinedError(w, err)
		}
		return
	}

	respJson, err := json.Marshal(t)
	if err != nil {
		helper.WriteUndefinedError(w, fmt.Errorf("could not convert response to json: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}

func NewTeamMux(handler *TeamHandler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/add", http.HandlerFunc(handler.add))
	mux.Handle("/get", http.HandlerFunc(handler.get))
	return mux
}
