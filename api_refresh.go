package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/33313/chirpy/internal/auth"
)

func (api *API) handleRefresh(w http.ResponseWriter, r *http.Request) {
	refresh, err := auth.GetBearerToken(r.Header.Get("Authorization"))
	if err != nil {
		log.Printf("Bearer token error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, ok := api.db.GetUserByToken(refresh)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := auth.CreateJWT(user.ID, api.jwtSecret)
	if err != nil {
		log.Printf("Error creating JWT: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
	if err != nil {
		handleJsonError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
