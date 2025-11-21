package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Ant-Shell/chirpy/internal/auth"
	"github.com/Ant-Shell/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request){
	type parameters struct {
			Body string `json:"body"`
	}

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or malformed Authorization header", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid or expired token", err)
		return
	}

	const maxChirpLength = 140
	body := strings.TrimSpace(params.Body)
	if body == "" {
		respondWithError(w, http.StatusBadRequest, "body is required", nil)
		return
	}

	if len(body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{"kerfuffle": {}, "sharbert": {}, "fornax": {}}
	cleaned := getCleanedBody(params.Body, badWords)

	dbChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleaned,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create chirp", err)
		return
	}

	chirp := Chirp{
		ID: dbChirp.ID,
		Body: dbChirp.Body,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID: dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word:= range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}