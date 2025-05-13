package main

import (
	"context"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	new_user, err := cfg.dbQueries.CreateUser(context.Background(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:        new_user.ID,
			CreatedAt: new_user.CreatedAt,
			UpdatedAt: new_user.UpdatedAt,
			Email:     new_user.Email,
		},
	})
	return

}
