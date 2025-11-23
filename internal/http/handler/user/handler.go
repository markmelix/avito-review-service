package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"review/internal/http/helper"
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
		helper.WriteMethodNotAllowedError(w)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		helper.WriteUndefinedError(w, fmt.Errorf("could not read request body: %w", err))
		return
	}
	defer r.Body.Close()

	var b SetIsActiveDto
	if err := json.Unmarshal(body, &b); err != nil {
		helper.WriteUndefinedError(w, fmt.Errorf("could not unmarshall request body: %w", err))
		return
	}

	ctx := context.Background()

	if err := h.usecase.SetIsActiveUser(ctx, b.UserId, b.IsActive); err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			helper.WriteNotFoundError(w)
		default:
			helper.WriteUndefinedError(w, err)
		}
		return
	}

	u, err := h.usecase.GetUser(ctx, b.UserId)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			helper.WriteNotFoundError(w)
		default:
			helper.WriteUndefinedError(w, err)
		}
		return
	}

	resp := &UserDto{User: *u}

	respJson, err := json.Marshal(resp)
	if err != nil {
		helper.WriteUndefinedError(w, fmt.Errorf("could not convert response to json: %w", err))
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

	uId := r.URL.Query().Get("user_id")
	if uId == "" {
		helper.WriteUndefinedError(w, errors.New("user_id query param was not specified"))
		return
	}

	ctx := context.Background()
	prs, err := h.usecase.GetReviewAssignments(ctx, uId)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			helper.WriteNotFoundError(w)
		default:
			helper.WriteUndefinedError(w, err)
		}
		return
	}

	resp := &GetReviewDto{UserId: uId, PullReqs: prs}

	respJson, err := json.Marshal(resp)
	if err != nil {
		helper.WriteUndefinedError(w, fmt.Errorf("could not convert response to json: %w", err))
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
