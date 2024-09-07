package database

import (
	"fmt"
)

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (db *DB) CreateUser(body string) (User, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return User{}, fmt.Errorf("problem loading db, %v", err)
	}

	lastId := 0
	for key := range dbStructure.Users {
		if key > lastId {
			lastId = key
		}
	}

	nextId := lastId + 1

	newUser := User{
		Id:    nextId,
		Email: body,
	}

	dbStructure.Users[nextId] = newUser

	err = db.writeDB(dbStructure)

	if err != nil {
		return User{}, fmt.Errorf("problem writing to db, %v", err)
	}

	return newUser, nil
}
