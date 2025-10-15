package main

import (
	"encoding/json"
	"unicode/utf8"
	"net/http"
)

func handlerValidate(w http.ResponseWriter, r *http.Request){
    type parameters struct {
        Body string `json:"body"`
    }

		type errVal struct {
				Error string `json:"error"`
		}

		type returnVal struct {
				Valid bool `json:"valid"`
		}

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
				errBody := errVal{
					Error: "Something went wrong",
				}
				errMsg, _ := json.Marshal(errBody)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(errMsg)
				return
    }
		runeCount := utf8.RuneCountInString(params.Body)
		if runeCount > 140 {
			errBody := errVal{
				Error: "Chirp is too long",
			}
			errMsg, _ := json.Marshal(errBody)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errMsg)
			return
		}

		returnBody := returnVal{
			Valid: true,
		}
		res, _ := json.Marshal(returnBody)
		w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
		w.Write(res)
}