package main

import "net/http"

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	dbChirpList, err := cfg.db.ListChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed get chirps", err)
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

	respondWithJSON(w, http.StatusOK, chirpList)
}