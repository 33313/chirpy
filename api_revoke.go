package main

import (
	"log"
	"net/http"

	"github.com/myshkovsky/chirpy/internal/auth"
	"github.com/myshkovsky/chirpy/internal/database"
)

func (api *API) handleRevoke(w http.ResponseWriter, r *http.Request) {
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

	api.db.UpdateUser(user.ID, database.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
		Refresh:  "",
		Red:      user.Red,
	})

	w.WriteHeader(http.StatusNoContent)
}
