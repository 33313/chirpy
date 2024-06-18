package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (api *fsAPI) handlePostUser(w http.ResponseWriter, r *http.Request) {
	type paramsPostUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := paramsPostUser{}
	decodeParams[paramsPostUser](w, r, &params)
	user, err := api.db.CreateUser(params.Email)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		handleJsonError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (api *fsAPI) handleGetUser(w http.ResponseWriter, r *http.Request) {
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
