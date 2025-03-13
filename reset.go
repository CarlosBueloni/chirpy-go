package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/plain; charset=utf-8")

	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Access Denied")
		return
	}

	cfg.dbQueries.Reset(r.Context())
}
