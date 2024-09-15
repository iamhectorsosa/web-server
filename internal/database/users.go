package database

import (
	"errors"
)

type User struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

var ErrUserAlreadyExists = errors.New("User already exists")
var ErrUserDoesNotExist = errors.New("User doesn't exist")
var ErrDatabaseLoad = errors.New("Error loading database")
var ErrDatabaseWrite = errors.New("Error writing to database")
var ErrPasswordMismatch = errors.New("Password doesn't match")

func (db *DB) CreateUser(email, passwordHash string) (User, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return User{}, ErrDatabaseLoad
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return User{}, ErrUserAlreadyExists
		}
	}

	lastId := 0
	for key := range dbStructure.Users {
		if key > lastId {
			lastId = key
		}
	}

	nextId := lastId + 1

	newUser := User{
		Id:           nextId,
		Email:        email,
		PasswordHash: passwordHash,
	}

	dbStructure.Users[nextId] = newUser

	err = db.writeDB(dbStructure)

	if err != nil {
		return User{}, ErrDatabaseWrite
	}

	return newUser, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return User{}, ErrDatabaseLoad
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		} else {
			return User{}, ErrUserDoesNotExist
		}

	}

	return User{}, ErrUserDoesNotExist
}

func (db *DB) UpdateUserById(userId int, email, passwordHash string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, ErrDatabaseLoad
	}

	_, ok := dbStructure.Users[userId]

	if !ok {
		return User{}, ErrUserDoesNotExist
	}

	newUser := User{
		Id:           userId,
		Email:        email,
		PasswordHash: passwordHash,
	}

	dbStructure.Users[userId] = newUser

	err = db.writeDB(dbStructure)

	if err != nil {
		return User{}, ErrDatabaseWrite
	}

	return newUser, nil
}
