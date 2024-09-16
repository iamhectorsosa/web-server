package database

import (
	"errors"
	"time"
)

type RefreshToken struct {
	UserId    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

var ErrRefreshTokenDoesNotExist = errors.New("Refresh token doesn't exist")

func (db *DB) CreateRefreshToken(userId int, token string, expiresAt time.Time) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return ErrDatabaseLoad
	}

	refreshToken := RefreshToken{
		UserId:    userId,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	dbStructure.RefreshTokens[token] = refreshToken

	err = db.writeDB(dbStructure)
	if err != nil {
		return ErrDatabaseWrite
	}

	return nil
}

func (db *DB) DeleteRefreshToken(token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return ErrDatabaseLoad
	}

	delete(dbStructure.RefreshTokens, token)

	err = db.writeDB(dbStructure)
	if err != nil {
		return ErrDatabaseWrite
	}

	return nil
}

func (db *DB) GetUserAndRefreshTokenByRefreshToken(token string) (User, RefreshToken, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, RefreshToken{}, ErrDatabaseLoad
	}

	refreshToken, ok := dbStructure.RefreshTokens[token]
	if !ok {
		return User{}, RefreshToken{}, ErrRefreshTokenDoesNotExist
	}

	user, err := db.GetUserById(refreshToken.UserId)
	if err != nil {
		return User{}, RefreshToken{}, err
	}

	return user, refreshToken, nil
}
