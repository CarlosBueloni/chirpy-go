package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	type returnVals struct {
		Cleaned_Body string `json:"cleaned_body"`
	}
	respondWithJSON(w, http.StatusOK, returnVals{
		Cleaned_Body: cleanText(params.Body),
	})
	return
}

func cleanText(text string) string {
	bad_words := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	words := strings.Split(text, " ")
	var clean_words []string
	for _, word := range words {
		_, ok := bad_words[strings.ToLower(word)]
		if ok {
			word = "****"
		}
		clean_words = append(clean_words, word)
	}
	return strings.Join(clean_words, " ")
}
