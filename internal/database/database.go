package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps        map[int]Chirp           `json:"chirps"`
	Users         map[int]User            `json:"users"`
	RefreshTokens map[string]RefreshToken `json:"refresh_tokens"`
}

var ErrDatabaseLoad = errors.New("Error loading database")
var ErrDatabaseWrite = errors.New("Error writing to database")

func NewDB(path string, debug bool) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}

	err := db.ensureDB(debug)
	return db, err
}

func (db *DB) ensureDB(debug bool) error {
	openFileFlags := os.O_RDWR | os.O_CREATE
	if debug {
		openFileFlags |= os.O_TRUNC
	}

	file, err := os.OpenFile(db.path, openFileFlags, 0666)

	if err != nil {
		return fmt.Errorf("problem opening %s, %v", db.path, err)
	}

	defer file.Close()

	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("problem getting file info from the file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 || debug {
		err = json.NewEncoder(file).Encode(DBStructure{
			Chirps:        map[int]Chirp{},
			Users:         map[int]User{},
			RefreshTokens: map[string]RefreshToken{},
		})

		if err != nil {
			return fmt.Errorf("problem encoding %s, %v", db.path, err)
		}
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
