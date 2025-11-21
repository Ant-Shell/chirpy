package main

import (
	"net/http"
	"strings"
)


func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "missing authorization header", nil)
		return
	}

	const prefix = "Bearer "
	token, ok := strings.CutPrefix(authHeader, prefix)
	if !ok || strings.TrimSpace(token) == "" {
		respondWithError(w, http.StatusUnauthorized, "invalid authorization header", nil)
		return
	}

	refreshToken := strings.TrimSpace(token)
	_, err := cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not revoke token", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}