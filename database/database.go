package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func (db *DB) ensureDB() error {
	file, err := os.OpenFile(db.path, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return fmt.Errorf("problem opening %s, %v", db.path, err)
	}

	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("problem getting file info from the file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte(`{"chirps":{}}`))
	}

	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	file, err := os.OpenFile(db.path, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return DBStructure{}, fmt.Errorf("problem opening %s, %v", db.path, err)
	}

	var dbStructure DBStructure

	err = json.NewDecoder(file).Decode(&dbStructure)

	if err != nil {
		return DBStructure{}, fmt.Errorf("problem parsing db, %v", err)
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	file, err := os.OpenFile(db.path, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return fmt.Errorf("problem opening %s, %v", db.path, err)
	}

	err = json.NewEncoder(file).Encode(dbStructure)

	if err != nil {
		return fmt.Errorf("problem encoding %s, %v", db.path, err)
	}

	return nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return Chirp{}, fmt.Errorf("problem loading db, %v", err)
	}

	lastId := 0
	for key := range dbStructure.Chirps {
		if key > lastId {
			lastId = key
		}
	}

	nextId := lastId + 1

	newChirp := Chirp{
		Id:   nextId,
		Body: body,
	}

	dbStructure.Chirps[nextId] = newChirp

	err = db.writeDB(dbStructure)

	if err != nil {
		return Chirp{}, fmt.Errorf("problem writing to db, %v", err)
	}

	return newChirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return nil, fmt.Errorf("problem loading db, %v", err)
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))

	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirpById(chirpId int) (Chirp, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return Chirp{}, fmt.Errorf("problem loading db, %v", err)
	}

	chirp, ok := dbStructure.Chirps[chirpId]

	if !ok {
		return Chirp{}, fmt.Errorf("chirp not found")
	}

	return chirp, nil
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}

	err := db.ensureDB()

	if err != nil {
		return nil, fmt.Errorf("problem parsing file, %v", err)
	}

	return db, nil
}
