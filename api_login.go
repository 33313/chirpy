package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/myshkovsky/chirpy/internal/auth"
	"github.com/myshkovsky/chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type UserSanitized struct {
	ID      int    `json:"id"`
	Email   string `json:"email"`
	Token   string `json:"token"`
	Refresh string `json:"refresh_token"`
	Red     bool   `json:"is_chirpy_red"`
}

func (api *API) handleLogin(w http.ResponseWriter, r *http.Request) {
	type requestParams struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
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

	token, err := auth.CreateJWT(user.ID, api.jwtSecret)
	refreshToken := auth.CreateRefreshToken()

	res, err := json.Marshal(UserSanitized{
		ID:      user.ID,
		Email:   user.Email,
		Token:   token,
		Refresh: refreshToken,
		Red:     user.Red,
	})
	if err != nil {
		handleJsonError(w, err)
		return
	}

	_, err = api.db.UpdateUser(user.ID, database.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
		Refresh:  refreshToken,
		Red:      user.Red,
	})
	if err != nil {
		log.Println("Error updating user: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
