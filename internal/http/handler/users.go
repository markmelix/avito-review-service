package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"review/internal/service/user"
)

type UserHandler struct {
	usecase *user.UserService
}

func NewUserHandler(usecase *user.UserService) *UserHandler {
	return &UserHandler{usecase}
}

func (h *UserHandler) setIsActive(w http.ResponseWriter, r *http.Request) {
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

func (h *UserHandler) getReview(w http.ResponseWriter, r *http.Request) {
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

func NewUserMux(handler *UserHandler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/setIsActive", http.HandlerFunc(handler.setIsActive))
	mux.Handle("/getReview", http.HandlerFunc(handler.getReview))
	return mux
}
