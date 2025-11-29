package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/Ant-Shell/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUpgradeChirpyRed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "an API key is required", err)
	}

	if key != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "invalid API key", nil)
	}

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	parsedID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to parse ID", err)
		return
	}

	_, err = cfg.db.UpdateUserChirpyRedMembership(r.Context(), parsedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "user not found", err)
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, "issue detected while trying to locate user", err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
