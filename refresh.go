package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/Ant-Shell/chirpy/internal/auth"
)


func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
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
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "user not found", nil)
		return
	}

	newAccessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create token", nil)
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{Token: newAccessToken})
}