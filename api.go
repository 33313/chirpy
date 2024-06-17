package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/myshkovsky/chirpy/internal/database"
)

type fsAPI struct {
	hits int
	db   *database.DB
}

func (api *fsAPI) mwMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.hits++
		next.ServeHTTP(w, r)
	})
}

func (api *fsAPI) handleDisplayMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
    <html>
        <body>
            <h1>Welcome, Chirpy Admin</h1>
            <p>Chirpy has been visited %d times!</p>
        </body>
    </html>
`, api.hits)))
}

func (api *fsAPI) handleResetMetrics(w http.ResponseWriter, r *http.Request) {
	api.hits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func (api *fsAPI) handlePostChirp(w http.ResponseWriter, r *http.Request) {
	type jsonParams struct {
		Body string `json:"body"`
	}
	type fail struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	params := jsonParams{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

func (api *fsAPI) handleGetChirps(w http.ResponseWriter, r *http.Request) {
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

func (api *fsAPI) handleGetChirp(w http.ResponseWriter, r *http.Request) {
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

func handleJsonError(w http.ResponseWriter, err error) {
		log.Printf("Error marshalling json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
