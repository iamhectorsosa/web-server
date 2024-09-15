package main

import (
	"encoding/json"
	"net/http"

	"github.com/iamhectorsosa/web-server/internal/auth"
)

func (api *apiConfig) postLogin(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "JSON decoding failed")
		return
	}

	user, err := api.DB.GetUserByEmail(payload.Email)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email and/or password combination")
		return
	}

	err = auth.CheckHashPassword(payload.Password, user.PasswordHash)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email and/or password combination")
		return
	}

	token, err := auth.CreateJWT(user.Id, api.jwtSecret, payload.ExpiresInSeconds)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "JWT Token Creation failed")
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
	}{
		Id:    user.Id,
		Email: user.Email,
		Token: token,
	})
}
