package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"review/internal/service/team"
)

type TeamHandler struct {
	usecase *team.TeamService
}

func NewTeamHandler(usecase *team.TeamService) *TeamHandler {
	return &TeamHandler{usecase}
}

func (h *TeamHandler) add(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp := make(map[string]any)
	respJson, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("error converting response to json: %v", err), http.StatusBadRequest)
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

	resp := make(map[string]any)
	respJson, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("error converting response to json: %v", err), http.StatusBadRequest)
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
