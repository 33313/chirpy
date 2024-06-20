package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

func NewDB(path string) *DB {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	db.ensureDB()
	return db
}

func (db *DB) createDB() {
	dbStruct := DBStructure{
		Chirps: make(map[int]Chirp),
		Users:  make(map[int]User),
	}
	db.writeDB(dbStruct)
}

func (db *DB) ensureDB() {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		db.createDB()
	}
}

func (db *DB) loadDB() DBStructure {
	db.mux.RLock()
	defer db.mux.RUnlock()
	db.ensureDB()
	dbStruct := DBStructure{}
	data, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		log.Fatalf("ERROR: Database file not found: %s", err)
	}
	err = json.Unmarshal(data, &dbStruct)
	if err != nil {
		log.Fatalf("ERROR: Failed to unmarshal database: %s", err)
	}
	return dbStruct
}

func (db *DB) writeDB(dbStruct DBStructure) {
	data, err := json.Marshal(dbStruct)
	if err != nil {
		log.Fatalf("ERROR: Failed to marshal DB for write: %s", err)
	}
	err = os.WriteFile(db.path, data, 0600)
	if err != nil {
		log.Fatalf("ERROR: Failed to write DB: %s", err)
	}
}
