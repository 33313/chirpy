package database

import (
	"regexp"
)

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

func (db *DB) GetChirps() []Chirp {
	data := db.loadDB()
	var chirps []Chirp
	for _, v := range data.Chirps {
		chirps = append(chirps, v)
	}
	return chirps
}

func (db *DB) GetChirp(id int) (Chirp, bool) {
	data := db.loadDB()
	chirp, ok := data.Chirps[id]
	return chirp, ok
}

func (db *DB) CreateChirp(body string, authorID int) (Chirp, error) {
	data := db.loadDB()
	nextID := len(data.Chirps) + 1
	chirp := Chirp{
		ID:       nextID,
		Body:     cleanChirp(body),
		AuthorID: authorID,
	}
	data.Chirps[nextID] = chirp
	db.writeDB(data)
	return chirp, nil
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
	re := regexp.MustCompile(`(?i)` + badword)
	replacement := "****"
	return re.ReplaceAllString(msg, replacement)
}
