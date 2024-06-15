package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type fsAPI struct {
	hits int
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

func (api *fsAPI) handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type jsonParams struct {
		Body string `json:"body"`
	}
	type success struct {
		Cleaned string `json:"cleaned_body"`
	}
	type fail struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	params := jsonParams{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > 140 {
		data, err := json.Marshal(fail{
			Error: "Chirp is too long",
		})
		if err != nil {
			log.Printf("Error marshalling json: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(data)
		return
	}
	res, err := json.Marshal(success{
        Cleaned: cleanChirp(params.Body),
	})
	if err != nil {
		log.Printf("Error marshalling json: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(res)
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func cleanChirp(msg string) string {
    badwords := [3]string{"kerfuffle", "sharbert", "fornax"}
    clean := msg
    for _, word := range badwords {
        clean = cleanWord(clean, word)
    }
    return clean
}

func cleanWord(msg string, badword string) string {
    re := regexp.MustCompile(`(?i)`+badword)
    replacement := "****"
    return re.ReplaceAllString(msg, replacement)
}
