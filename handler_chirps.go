package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/CarlosBueloni/chirpy-go/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type response struct {
		Chirp
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

	new_chirp, err := cfg.dbQueries.CreateChirp(context.Background(), database.CreateChirpParams{
		Body:   cleanText(params.Body),
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Chirp: Chirp{
			ID:        new_chirp.ID,
			CreatedAt: new_chirp.CreatedAt,
			UpdatedAt: new_chirp.UpdatedAt,
			Body:      new_chirp.Body,
			UserID:    new_chirp.UserID,
		},
	})
	return
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirp_id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}
	chirp, err := cfg.dbQueries.GetChirp(r.Context(), chirp_id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}
	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
	return
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Chirps []Chirp
	}
	res, err := cfg.dbQueries.GetChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	chirps := []Chirp{}
	for _, chirp := range res {
		chirps = append(chirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, chirps)
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
