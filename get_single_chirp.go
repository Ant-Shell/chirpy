package main

import (
	"net/http"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetSingleChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	parsedID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to parse ID", err)
	}

	singleChirp, err := cfg.db.GetChirp(r.Context(), parsedID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "unable to locate chirp", err)
	}

	foundChirp := Chirp{
		ID: singleChirp.ID,
		Body: singleChirp.Body,
		CreatedAt: singleChirp.CreatedAt,
		UpdatedAt: singleChirp.UpdatedAt,
		UserID: singleChirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, foundChirp)
}