package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/myshkovsky/chirpy/internal/database"
)

type fsAPI struct {
	hits int
	db   *database.DB
}

type fail struct {
	Error string `json:"error"`
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

func decodeParams[T any](w http.ResponseWriter, r *http.Request, data *T) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
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
