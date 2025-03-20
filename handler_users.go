package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/carlosbueloni/chirpy-go/internal/auth"
	"github.com/carlosbueloni/chirpy-go/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	dat, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't read request")
		return
	}

	params := parameters{}
	err = json.Unmarshal(dat, &params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't unmarshal json")
		return
	}

	psw, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password")
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		Email:          params.Email,
		HashedPassword: psw,
	})
	if err != nil {
		log.Printf("error:%v", err)
		respondWithError(w, http.StatusInternalServerError, "error creating user")
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	dat, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't read request")
		return
	}

	params := parameters{}
	err = json.Unmarshal(dat, &params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't unmarshal json")
		return
	}

	user, err := cfg.dbQueries.GetUserFromEmail(r.Context(), params.Email)
	if err != nil {
		log.Printf("error:%v", err)
		respondWithError(w, http.StatusInternalServerError, "error fetching user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}
