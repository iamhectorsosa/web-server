package database

import (
	"fmt"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
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
