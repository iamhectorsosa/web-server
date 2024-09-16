package database

import (
	"errors"
)

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

var ErrChirpDoesNotExist = errors.New("Chirp doesn't exist")

func (db *DB) CreateChirp(body string, authorId int) (Chirp, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return Chirp{}, ErrDatabaseLoad
	}

	lastId := 0
	for key := range dbStructure.Chirps {
		if key > lastId {
			lastId = key
		}
	}

	nextId := lastId + 1

	newChirp := Chirp{
		Id:       nextId,
		Body:     body,
		AuthorId: authorId,
	}

	dbStructure.Chirps[nextId] = newChirp

	err = db.writeDB(dbStructure)

	if err != nil {
		return Chirp{}, ErrDatabaseWrite
	}

	return newChirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return nil, ErrDatabaseLoad
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
		return Chirp{}, ErrDatabaseLoad
	}

	chirp, ok := dbStructure.Chirps[chirpId]

	if !ok {
		return Chirp{}, ErrChirpDoesNotExist
	}

	return chirp, nil
}

func (db *DB) DeleteChirpById(chirpId int) error {
	dbStructure, err := db.loadDB()

	if err != nil {
		return ErrDatabaseLoad
	}

	delete(dbStructure.Chirps, chirpId)

	err = db.writeDB(dbStructure)
	if err != nil {
		return ErrDatabaseWrite
	}

	return nil
}
