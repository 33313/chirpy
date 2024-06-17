package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"regexp"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	chirps := db.GetChirps()
	lastID := len(chirps) + 1
	chirp := Chirp{
		ID:   lastID,
		Body: cleanChirp(body),
	}
	chirps = append(chirps, chirp)
	chirpsToSave := make(map[int]Chirp)
	for _, v := range chirps {
		chirpsToSave[v.ID] = v
	}
	db.writeDB(DBStructure{
		Chirps: chirpsToSave,
	})
	return chirp, nil
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

func (db *DB) createDB() error {
	dbStruct := DBStructure{
		Chirps: map[int]Chirp{},
	}
	return db.writeDB(dbStruct)
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	err := db.ensureDB()

	dbStruct := DBStructure{}
	data, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStruct, err
	}

	err = json.Unmarshal(data, &dbStruct)
	if err != nil {
		return dbStruct, err
	}
	return dbStruct, nil
}

func (db *DB) writeDB(dbStruct DBStructure) error {
	data, err := json.Marshal(dbStruct)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, 0600)
	if err != nil {
		return err
	}
	return nil
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
