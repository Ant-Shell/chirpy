package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	dbChirpList, err := cfg.db.ListChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed get chirps", err)
		return
	}

	chirpList := []Chirp{}
	for _, chirp := range dbChirpList {
		updatedChirp := Chirp{
			ID: chirp.ID,
			Body: chirp.Body,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			UserID: chirp.UserID,
		}
		chirpList = append(chirpList, updatedChirp)
	}

	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString == "" {
		respondWithJSON(w, http.StatusOK, chirpList)
		return
	}

	parsedID, err := uuid.Parse(authorIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid author_id", err)
		return
	}

	targetedChirpList := []Chirp{}
	for _, targetedChirp := range chirpList {
		if targetedChirp.UserID == parsedID {
			targetedChirpList = append(targetedChirpList, targetedChirp)
		}
	}

	respondWithJSON(w, http.StatusOK, targetedChirpList)
}