package pr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"review/internal/http/helper"
	"review/internal/service/pr"
)

type PullReqHandler struct {
	usecase *pr.PullReqService
}

func NewPullReqHandler(usecase *pr.PullReqService) *PullReqHandler {
	return &PullReqHandler{usecase}
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

	resp := CreateRespDtoFromPullReq(prObj)
	respJson, err := json.Marshal(resp)
	if err != nil {
		helper.WriteUndefinedError(w, fmt.Errorf("could not convert response to json: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}

func (h *PullReqHandler) merge(w http.ResponseWriter, r *http.Request) {
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

	var b PullReqMergeReqDto
	if err := json.Unmarshal(body, &b); err != nil {
		helper.WriteUndefinedError(w, fmt.Errorf("could not unmarshall request body: %w", err))
		return
	}

	ctx := context.Background()

	prObj, err := h.usecase.MergePullReq(ctx, b.PullReqId)
	if err != nil {
		switch {
		case errors.Is(err, pr.ErrPullReqNotFound):
			helper.WriteNotFoundError(w)
		default:
			helper.WriteUndefinedError(w, err)
		}
		return
	}

	resp := MergeRespDtoFromPullReq(prObj)
	respJson, err := json.Marshal(resp)
	if err != nil {
		helper.WriteUndefinedError(w, fmt.Errorf("could not convert response to json: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}

func (h *PullReqHandler) reassign(w http.ResponseWriter, r *http.Request) {
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

	var b PullReqReassignReqDto
	if err := json.Unmarshal(body, &b); err != nil {
		helper.WriteUndefinedError(w, fmt.Errorf("could not unmarshall request body: %w", err))
		return
	}

	ctx := context.Background()

	prObj, replacedById, err := h.usecase.ReassignReviewer(ctx, b.PullReqId, b.OldReviewerId)

	if err != nil {
		switch {
		case errors.Is(err, pr.ErrPullReqNotFound) || errors.Is(err, pr.ErrAuthorNotFound):
			helper.WriteNotFoundError(w)
		case errors.Is(err, pr.ErrReassignOnMerged):
			helper.WriteError(w, "PR_MERGED", "cannot reassign on merged PR", http.StatusConflict)
		case errors.Is(err, pr.ErrReviewerNotAssigned):
			helper.WriteError(w, "NOT_ASSIGNED", "reviewer is not assigned to this PR", http.StatusConflict)
		case errors.Is(err, pr.ErrNoReassignCandidate):
			helper.WriteError(w, "NO_CANDIDATE", "no active replacement candidate in team", http.StatusConflict)
		default:
			helper.WriteUndefinedError(w, err)
		}
		return
	}

	resp := ReassignRespFromPullReq(prObj, replacedById)
	respJson, err := json.Marshal(resp)
	if err != nil {
		helper.WriteUndefinedError(w, fmt.Errorf("could not convert response to json: %w", err))
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
