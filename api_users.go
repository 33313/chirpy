package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/myshkovsky/chirpy/internal/auth"
	"github.com/myshkovsky/chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (api *API) handlePostUser(w http.ResponseWriter, r *http.Request) {
	type postUserSanitized struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}
	type requestParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := requestParams{}
	decodeParams[requestParams](w, r, &params)

	user, err := api.db.CreateUser(params.Email, params.Password)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(postUserSanitized{
		ID:    user.ID,
		Email: user.Email,
	})
	if err != nil {
		handleJsonError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (api *API) handleGetUser(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Fatalf("Error converting string->int: %s", err)
	}

	user, ok := api.db.GetUser(n)
	if !ok {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 Not Found"))
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		handleJsonError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (api *API) handlePutUser(w http.ResponseWriter, r *http.Request) {
	type requestParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := requestParams{}
	decodeParams[requestParams](w, r, &params)

	token, err := auth.GetBearerToken(r.Header.Get("Authorization"))
	if err != nil {
		log.Printf("Error updating user: %s", err)
		w.Header().Set("Content-Type", "text/plain")
		switch err {
		case auth.ErrNoAuthHeader:
			w.WriteHeader(http.StatusBadRequest)
		case auth.ErrBadAuthHeader:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))
	}

	idStr, err := auth.ValidateJWT(token, api.jwtSecret)
	if err != nil {
		log.Printf("Error validating token: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Error parsing str->int: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	oldUser, ok := api.db.GetUser(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	pwd, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	api.db.UpdateUser(id, database.User{
		ID:       id,
		Email:    params.Email,
		Password: pwd,
        Refresh: oldUser.Refresh,
	})

	res, err := json.Marshal(database.User{
		ID:    id,
		Email: params.Email,
	})
	if err != nil {
		handleJsonError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
