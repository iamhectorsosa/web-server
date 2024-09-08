package main

import (
	"encoding/json"
	"net/http"
)

func (api *apiConfig) postUsers(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Email string `json:"email"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	user, err := api.DB.CreateUser(payload.Email)

	respondWithJSON(w, http.StatusCreated, user)
}
