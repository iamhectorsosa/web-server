package database

import (
	"errors"
)

type User struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
}

var ErrUserAlreadyExists = errors.New("User already exists")
var ErrUserDoesNotExist = errors.New("User doesn't exist")
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

func (db *DB) GetUserById(userId int) (User, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return User{}, ErrDatabaseLoad
	}

	for _, user := range dbStructure.Users {
		if user.Id == userId {
			return user, nil
		} else {
			return User{}, ErrUserDoesNotExist
		}

	}

	return User{}, ErrUserDoesNotExist
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return User{}, ErrDatabaseLoad
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, ErrUserDoesNotExist
}

func (db *DB) UpdateUserEmailPasswordById(userId int, email, passwordHash string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, ErrDatabaseLoad
	}

	user, ok := dbStructure.Users[userId]

	if !ok {
		return User{}, ErrUserDoesNotExist
	}

	updatedUser := User{
		Id:           user.Id,
		Email:        email,
		PasswordHash: passwordHash,
	}

	dbStructure.Users[userId] = updatedUser

	err = db.writeDB(dbStructure)

	if err != nil {
		return User{}, ErrDatabaseWrite
	}

	return updatedUser, nil
}

func (db *DB) UpgradeUserToRedByUserId(userId int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return ErrDatabaseLoad
	}

	user, ok := dbStructure.Users[userId]

	if !ok {
		return ErrUserDoesNotExist
	}

	if user.IsChirpyRed {
		return nil
	}

	upgradedUser := User{
		Id:           user.Id,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		IsChirpyRed:  true,
	}

	dbStructure.Users[userId] = upgradedUser

	err = db.writeDB(dbStructure)

	if err != nil {
		return ErrDatabaseWrite
	}

	return nil
}
