package main

import (
	"encoding/json"
	"net/http"

	"github.com/iamhectorsosa/web-server/internal/auth"
)

func (api *apiConfig) postUsers(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "JSON decoding failed")
		return
	}

	passwordHash, err := auth.HashPassword(payload.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Password hashing failed")
		return
	}

	user, err := api.DB.CreateUser(payload.Email, passwordHash)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "User creation failed")
		return
	}

	respondWithJSON(w, http.StatusCreated, struct {
		Id          int    `json:"id"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}{
		Id:          user.Id,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (api *apiConfig) putUsers(w http.ResponseWriter, r *http.Request) {

	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthenticated request")
		return
	}

	userId, err := auth.ValidateJWT(authToken, api.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthenticated request")
		return
	}

	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "JSON decoding failed")
		return
	}

	passwordHash, err := auth.HashPassword(payload.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Password hashing failed")
		return
	}

	user, err := api.DB.UpdateUserEmailPasswordById(userId, payload.Email, passwordHash)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating User")
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Id          int    `json:"id"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}{
		Id:          user.Id,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
