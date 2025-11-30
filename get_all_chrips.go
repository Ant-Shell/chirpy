package main

import (
	"net/http"
	"sort"

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
	if authorIDString != "" {
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
		chirpList = targetedChirpList
	}

	sort_order := r.URL.Query().Get("sort")
	if sort_order == "desc" {
		sort.Slice(chirpList, func(i, j int) bool {
			return chirpList[i].CreatedAt.After(chirpList[j].CreatedAt)
		})
	} else {
		sort.Slice(chirpList, func(i, j int) bool {
			return chirpList[i].CreatedAt.Before(chirpList[j].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, chirpList)
}