package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/myshkovsky/chirpy/internal/auth"
)

func (api *API) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header.Get("Authorization"))
	if err != nil {
		log.Printf("Bearer token error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	idStr, err := auth.ValidateJWT(token, api.jwtSecret)
	if err != nil {
		log.Printf("Error validating JWT: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Error converting str->int: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	n, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Fatalf("Error converting string->int: %s", err)
	}

	chirp, ok := api.db.GetChirp(n)
	if !ok {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 Not Found"))
		return
	}

	if chirp.AuthorID != userID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	api.db.DeleteChirp(n)
	w.WriteHeader(http.StatusNoContent)
}
