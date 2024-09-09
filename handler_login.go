package main

import (
	"encoding/json"
	"net/http"
)

func (api *apiConfig) postLogin(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	user, err := api.DB.GetUserByEmailPassword(payload.Email, payload.Password)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, UserResponse{
		Id:    user.Id,
		Email: user.Email,
	})
}
