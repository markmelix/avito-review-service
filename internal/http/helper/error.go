package helper

import (
	"encoding/json"
	"net/http"
)

type ErrorRespBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResp struct {
	Error ErrorRespBody `json:"error"`
}

func WriteError(w http.ResponseWriter, code, message string, httpStatus int) {
	body := ErrorResp{Error: ErrorRespBody{Code: code, Message: message}}

	resp, err := json.Marshal(body)
	if err != nil {
		http.Error(w, "failed to marshall error response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write(resp)
}

func WriteNotFoundError(w http.ResponseWriter) {
	WriteError(w, "NOT_FOUND", "resource not found", http.StatusNotFound)
}

func WriteUndefinedError(w http.ResponseWriter, err error) {
	WriteError(w, "UNDEFINED", err.Error(), http.StatusInternalServerError)
}

func WriteMethodNotAllowedError(w http.ResponseWriter) {
	WriteError(w, "METHOD_NOT_ALLOWED", "method not allowed", http.StatusMethodNotAllowed)
}
