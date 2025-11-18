package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Ant-Shell/chirpy/internal/auth"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
			Email string `json:"email"`
			Password string `json:"password"`
			ExpiresInSeconds *int64 `json:"expires_in_seconds"`
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

	max := int64(3600)
	exp := max
	if p.ExpiresInSeconds != nil {
		v := min(*p.ExpiresInSeconds, max)
		if v < 1 {
			v = 1
		}
		exp = v
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

	token, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, time.Duration(exp)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create token.", nil)
		return
	}

	user := loginResponse{
		ID: dbUser.ID,
		Email: dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Token: token,
	}

	respondWithJSON(w, http.StatusOK, user)
}