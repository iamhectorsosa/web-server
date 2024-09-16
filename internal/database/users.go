package database

import (
	"errors"
	"time"
)

type User struct {
	Id                     int       `json:"id"`
	Email                  string    `json:"email"`
	PasswordHash           string    `json:"password_hash"`
	RefreshToken           string    `json:"refresh_token"`
	RefreshTokenExpiration time.Time `json:"refresh_token_expiration"`
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
		Id:                     user.Id,
		Email:                  email,
		PasswordHash:           passwordHash,
		RefreshToken:           user.RefreshToken,
		RefreshTokenExpiration: user.RefreshTokenExpiration,
	}

	dbStructure.Users[userId] = updatedUser

	err = db.writeDB(dbStructure)

	if err != nil {
		return User{}, ErrDatabaseWrite
	}

	return updatedUser, nil
}

func (db *DB) UpdateUserRefreshTokenById(userId int, refreshToken string, refreshTokenExpiration time.Time) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, ErrDatabaseLoad
	}

	user, ok := dbStructure.Users[userId]

	if !ok {
		return User{}, ErrUserDoesNotExist
	}

	updatedUser := User{
		Id:                     userId,
		Email:                  user.Email,
		PasswordHash:           user.PasswordHash,
		RefreshToken:           refreshToken,
		RefreshTokenExpiration: refreshTokenExpiration,
	}

	dbStructure.Users[userId] = updatedUser

	err = db.writeDB(dbStructure)

	if err != nil {
		return User{}, ErrDatabaseWrite
	}

	return updatedUser, nil
}

func (db *DB) GetUserByRefreshToken(refreshToken string) (User, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return User{}, ErrDatabaseLoad
	}

	for _, user := range dbStructure.Users {
		if user.RefreshToken == refreshToken {
			return user, nil
		} else {
			return User{}, ErrUserDoesNotExist
		}

	}

	return User{}, ErrUserDoesNotExist
}

func (db *DB) DeleteRefreshTokenByRefreshToken(refreshToken string) error {
	dbStructure, err := db.loadDB()

	if err != nil {
		return ErrDatabaseLoad
	}

	for _, user := range dbStructure.Users {
		if user.RefreshToken == refreshToken {
			dbStructure.Users[user.Id] = User{
				Id:                     user.Id,
				Email:                  user.Email,
				PasswordHash:           user.PasswordHash,
				RefreshToken:           "",
				RefreshTokenExpiration: time.Now().UTC(),
			}

			err = db.writeDB(dbStructure)
			if err != nil {
				return ErrDatabaseWrite
			}

			return nil
		} else {
			return ErrUserDoesNotExist
		}

	}

	return ErrUserDoesNotExist
}
