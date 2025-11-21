package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Ant-Shell/chirpy/internal/auth"
	"github.com/Ant-Shell/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
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

	hash, err := auth.HashPassword(p.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not hash password", err)
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: p.Email,
		HashedPassword: hash,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create user", err)
		return
	}

		user := userResponse{
		ID: dbUser.ID,
		Email: dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}

	respondWithJSON(w, http.StatusCreated, user)
}