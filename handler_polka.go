package main

import (
	"encoding/json"
	"net/http"

	"github.com/iamhectorsosa/web-server/internal/auth"
	database "github.com/iamhectorsosa/web-server/internal/database"
)

var UserUpgradedEvent = "user.upgraded"

func (api *apiConfig) postUserUpgrade(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil || apiKey != api.polkaApiKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	payload := struct {
		Event string `json:"event"`
		Data  struct {
			UserId int `json:"user_id"`
		} `json:"data"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "JSON decoding failed")
		return
	}

	if payload.Event != UserUpgradedEvent {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = api.DB.UpgradeUserToRedByUserId(payload.Data.UserId)

	if err != nil && err == database.ErrUserDoesNotExist {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
