package main

import (
	"encoding/json"
	"unicode/utf8"
	"net/http"
	"strings"
	"slices"
)

func handlerValidate(w http.ResponseWriter, r *http.Request){
    type parameters struct {
        Body string `json:"body"`
    }

		type errVal struct {
				Error string `json:"error"`
		}

		type returnVal struct {
				Cleaned_Body string `json:"cleaned_body"`
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

		badWords := []string{"kerfuffle", "sharbert", "fornax"}
		body := strings.Split(params.Body, " ")
		cleaned_body := []string{}

		for i := range body {
			if slices.Contains(badWords, strings.ToLower(body[i])) {
				cleaned_body = append(cleaned_body, "****")
			} else {
				cleaned_body = append(cleaned_body, body[i])
			}
		}

		returnBody := returnVal{
			Cleaned_Body: strings.Join(cleaned_body, " "),
		}

		res, _ := json.Marshal(returnBody)
		w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
		w.Write(res)
}
