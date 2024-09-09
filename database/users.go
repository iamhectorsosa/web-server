package database

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

var ErrUserAlreadyExists = errors.New("User already exists")
var ErrUserDoesNotExist = errors.New("User doesn't exist")

func (db *DB) CreateUser(email, password string) (User, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return User{}, fmt.Errorf("problem loading db, %v", err)
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

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return User{}, fmt.Errorf("problem hashing password, %v", err)
	}

	newUser := User{
		Id:           nextId,
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	dbStructure.Users[nextId] = newUser

	err = db.writeDB(dbStructure)

	if err != nil {
		return User{}, fmt.Errorf("problem writing to db, %v", err)
	}

	return newUser, nil
}

func (db *DB) GetUserByEmailPassword(email, password string) (User, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return User{}, fmt.Errorf("problem loading db, %v", err)
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
			if err != nil {
				return User{}, fmt.Errorf("email/password combination didn't match, %v", err)
			}

			return user, nil
		}
	}

	return User{}, ErrUserDoesNotExist
}
