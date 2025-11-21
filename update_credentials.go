package main

import (
	"encoding/json"
	"net/http"

	"github.com/Ant-Shell/chirpy/internal/auth"
	"github.com/Ant-Shell/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateCredentials(w http.ResponseWriter, r *http.Request ) {
	type parameters struct {
			Email string `json:"email"`
			Password string `json:"password"`
	}

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid or missing token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid or missing token", err)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not hash password", err)
		return
	}

	user, err := cfg.db.UpdateUserCredentials(r.Context(), database.UpdateUserCredentialsParams{
		ID: userID,
		Email: params.Email,
		HashedPassword: hash,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update credentials", err)
		return
	}

	updatedUser := userResponse{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	respondWithJSON(w, http.StatusOK, updatedUser)
}