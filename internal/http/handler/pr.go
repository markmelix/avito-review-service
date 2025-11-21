package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"review/internal/service"
)

type PullReqHandler struct {
	usecase *service.PullReqService
}

func NewPullReqHandler(usecase *service.PullReqService) *PullReqHandler {
	return &PullReqHandler{usecase}
}

func (h *PullReqHandler) create(w http.ResponseWriter, r *http.Request) {
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

func (h *PullReqHandler) merge(w http.ResponseWriter, r *http.Request) {
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

func (h *PullReqHandler) reassign(w http.ResponseWriter, r *http.Request) {
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

func NewPullReqMux(handler *PullReqHandler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/create", http.HandlerFunc(handler.create))
	mux.Handle("/merge", http.HandlerFunc(handler.merge))
	mux.Handle("/reassign", http.HandlerFunc(handler.reassign))
	return mux
}
