package handler

import (
	"net/http"
	"slices"
	"strings"
)

var originAllowlist = []string{"*"}
var methodAllowlist = []string{"GET", "POST", "DELETE", "OPTIONS"}

func isPreflight(r *http.Request) bool {
	return r.Method == "OPTIONS" &&
		r.Header.Get("Origin") != "" &&
		r.Header.Get("Access-Control-Request-Method") != ""
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPreflight(r) {
			origin := r.Header.Get("Origin")
			method := r.Header.Get("Access-Control-Request-Method")
			if (slices.Contains(originAllowlist, "*") || slices.Contains(originAllowlist, origin)) && slices.Contains(methodAllowlist, method) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(methodAllowlist, ", "))
			}
		} else {
			origin := r.Header.Get("Origin")
			if slices.Contains(originAllowlist, "*") || slices.Contains(originAllowlist, origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
		}
		w.Header().Add("Vary", "Origin")
		next.ServeHTTP(w, r)
	})
}
