package main

import (
	"net/http"

	"github.com/Ant-Shell/chirpy/internal/auth"
	"github.com/Ant-Shell/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	parsedID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or malformed Authorization header", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid or missing token", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), parsedID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "unable to locate chirp", err)
		return
	}

	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "user does not have permission to access this resource", nil)
		return
	}

	num, err := cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID: chirp.ID,
		UserID: userID,
	})
	if err != nil || num == 0 {
		respondWithError(w, http.StatusInternalServerError, "unable to delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}