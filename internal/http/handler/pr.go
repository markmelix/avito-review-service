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
	"review/internal/service/pr"
)

type PullReqHandler struct {
	usecase *pr.PullReqService
}

func NewPullReqHandler(usecase *pr.PullReqService) *PullReqHandler {
	return &PullReqHandler{usecase}
}

type PullReqCreateReqDto struct {
	PullReqId   string `json:"pull_request_id"`
	PullReqName string `json:"pull_request_name"`
	AuthorId    string `json:"author_id"`
}

type PullReqCreateRespBody struct {
	PullReqId   string               `json:"pull_request_id"`
	PullReqName string               `json:"pull_request_name"`
	AuthorId    string               `json:"author_id"`
	Status      domain.PullReqStatus `json:"status"`
	ReviewerIds []string             `json:"assigned_reviewers"`
}

type PullReqCreateRespDto struct {
	PullReq PullReqCreateRespBody `json:"pr"`
}

func usersToIds(users []domain.User) []string {
	s := make([]string, len(users))
	for i, u := range users {
		s[i] = u.Id
	}
	return s
}

func (h *PullReqHandler) create(w http.ResponseWriter, r *http.Request) {
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

	var b PullReqCreateReqDto
	if err := json.Unmarshal(body, &b); err != nil {
		helper.WriteUndefinedError(w, fmt.Errorf("could not unmarshall request body: %w", err))
		return
	}

	ctx := context.Background()

	if err := h.usecase.CreatePullReq(ctx, b.PullReqId, b.PullReqName, b.AuthorId); err != nil {
		switch {
		case errors.Is(err, pr.ErrAuthorNotFound) || errors.Is(err, pr.ErrPullReqNotFound):
			helper.WriteNotFoundError(w)
		case errors.Is(err, pr.ErrPullReqAlreadyExists):
			helper.WriteError(w, "PR_EXISTS", "PR id already exists", http.StatusConflict)
		default:
			helper.WriteUndefinedError(w, err)
		}
		return
	}

	prObj, err := h.usecase.GetPullReq(ctx, b.PullReqId)
	if err != nil {
		switch {
		case errors.Is(err, pr.ErrPullReqNotFound):
			helper.WriteNotFoundError(w)
		default:
			helper.WriteUndefinedError(w, err)
		}
		return
	}

	resp := &PullReqCreateRespDto{
		PullReq: PullReqCreateRespBody{
			PullReqId:   prObj.Id,
			PullReqName: prObj.Name,
			AuthorId:    prObj.AuthorId,
			Status:      prObj.Status,
			ReviewerIds: usersToIds(prObj.Reviewers),
		},
	}
	respJson, err := json.Marshal(resp)
	if err != nil {
		helper.WriteUndefinedError(w, fmt.Errorf("could not convert response to json: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}

func (h *PullReqHandler) merge(w http.ResponseWriter, r *http.Request) {

}

func (h *PullReqHandler) reassign(w http.ResponseWriter, r *http.Request) {

}

func NewPullReqMux(handler *PullReqHandler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/create", http.HandlerFunc(handler.create))
	mux.Handle("/merge", http.HandlerFunc(handler.merge))
	mux.Handle("/reassign", http.HandlerFunc(handler.reassign))
	return mux
}
