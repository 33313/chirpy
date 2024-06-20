package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UserSanitized struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (api *fsAPI) handleLogin(w http.ResponseWriter, r *http.Request) {
	type requestParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := requestParams{}
	decodeParams[requestParams](w, r, &params)

	user, ok := api.db.GetUserByEmail(params.Email)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := bcrypt.CompareHashAndPassword(user.Password, []byte(params.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	res, err := json.Marshal(UserSanitized{
		ID:    user.ID,
		Email: user.Email,
	})
	if err != nil {
		handleJsonError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
