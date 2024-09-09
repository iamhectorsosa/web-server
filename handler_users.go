package main

import (
	"encoding/json"
	"net/http"
)

type UserResponse struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

func (api *apiConfig) postUsers(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	user, err := api.DB.CreateUser(payload.Email, payload.Password)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusCreated, UserResponse{
		Id:    user.Id,
		Email: user.Email,
	})
}
