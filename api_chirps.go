package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (api *API) handlePostChirp(w http.ResponseWriter, r *http.Request) {
	type requestParams struct {
		Body string `json:"body"`
	}
	params := requestParams{}
	decodeParams[requestParams](w, r, &params)

	if len(params.Body) > 140 {
		data, err := json.Marshal(fail{
			Error: "Chirp is too long",
		})
		if err != nil {
			handleJsonError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}

	chirp, err := api.db.CreateChirp(params.Body)
	if err != nil {
		log.Printf("Error creating Chirp: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(chirp)
	if err != nil {
		handleJsonError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (api *API) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps := api.db.GetChirps()
	res, err := json.Marshal(chirps)
	if err != nil {
		handleJsonError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (api *API) handleGetChirp(w http.ResponseWriter, r *http.Request) {
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
	res, err := json.Marshal(chirp)
	if err != nil {
		handleJsonError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
