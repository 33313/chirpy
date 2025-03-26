package main

import (
	"errors"
	"net/http"

	"github.com/33313/chirpy/internal/auth"
	"github.com/33313/chirpy/internal/database"
)

func (api *API) handlePolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header.Get("Authorization"))
	if err != nil {
		switch err {
		case auth.ErrNoAuthHeader:
			w.WriteHeader(http.StatusUnauthorized)
		case auth.ErrBadAuthHeader:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if token != api.polka {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	type requestParams struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}
	params := requestParams{}
	decodeParams[requestParams](w, r, &params)

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = api.db.UpgradeUser(params.Data.UserID)
	if errors.Is(err, database.ErrUserNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
