package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/CarlosBueloni/chirpy-go/internal/auth"
	"github.com/CarlosBueloni/chirpy-go/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	if ok := auth.CheckPasswordHash(params.Password, user.HashedPassword); ok != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	expirationTime := time.Hour

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.secret,
		expirationTime,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
	}

	rt, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
	}

	refresh_token, err := cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     rt,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token:        accessToken,
		RefreshToken: refresh_token.Token,
	})
	return
}
func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	var nullTime sql.NullTime

	bearer_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Couldn't find JWT", err)
		return
	}

	refresh_token, err := cfg.dbQueries.GetRefreshToken(r.Context(), bearer_token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh Token Not Found", err)
		return
	}
	if time.Now().UTC().After(refresh_token.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "Credentials expired please login again", err)
		return
	}

	if refresh_token.RevokedAt != nullTime {
		respondWithError(w, http.StatusUnauthorized, "Credentials Revoked please login again", err)
		return
	}

	user, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), refresh_token.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching user", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating JWT", err)
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: token,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {

	bearer_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	refresh_token, err := cfg.dbQueries.GetRefreshToken(r.Context(), bearer_token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unathorized", err)
		return
	}
	cfg.dbQueries.UpdateRevokedAt(r.Context(), refresh_token.Token)
	respondWithJSON(w, http.StatusNoContent, nil)
}
