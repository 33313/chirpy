package database

import (
	"log"
	"regexp"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (db *DB) GetChirps() []Chirp {
	data, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
	}

	var chirps []Chirp
	for _, v := range data.Chirps {
		chirps = append(chirps, v)
	}

	return chirps
}

func (db *DB) GetChirp(id int) (Chirp, bool) {
	data, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
	}

	chirp, ok := data.Chirps[id]
	return chirp, ok
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		log.Fatalf("Error creating user: %s", err)
		return Chirp{}, err
	}
	nextID := len(dbStruct.Chirps) + 1
	chirp := Chirp{
		ID:   nextID,
		Body: cleanChirp(body),
	}
	dbStruct.Chirps[nextID] = chirp
	db.writeDB(dbStruct)
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
