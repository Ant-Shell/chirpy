package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Ant-Shell/chirpy/internal/auth"
	"github.com/Ant-Shell/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
			Email string `json:"email"`
			Password string `json:"password"`
	}

	var p parameters
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON", err)
		return
	}
	p.Email = strings.TrimSpace(p.Email)
	p.Password = strings.TrimSpace(p.Password)

	if p.Email == "" {
		respondWithError(w, http.StatusBadRequest, "email is required", nil)
		return
	}

	if p.Password == "" {
		respondWithError(w, http.StatusBadRequest, "password is required", nil)
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(r.Context(), p.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", nil)
		return
	}

	ok, err := auth.CheckPasswordHash(p.Password, dbUser.HashedPassword)
	if err != nil || !ok {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", nil)
		return
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create token.", nil)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create refresh token.", nil)
		return
	}

	expiresAt := time.Now().Add(60 * 24 * time.Hour)

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: dbUser.ID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
    respondWithError(w, http.StatusInternalServerError, "failed to store refresh token", err)
    return
}

	user := loginResponse{
		ID: dbUser.ID,
		Email: dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Token: token,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, http.StatusOK, user)
}