package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func (api *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := api.DB.GetChirps()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (api *apiConfig) getChirpById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	chirpId, err := strconv.Atoi(id)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID")
		return
	}

	chirps, err := api.DB.GetChirpById(chirpId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (api *apiConfig) postChirps(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	payload := struct {
		Body string `json:"body"`
	}{}

	err := decoder.Decode(&payload)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(payload.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	profaneWords := []string{
		"kerfuffle", "sharbert", "fornax",
	}

	cleanedBody := payload.Body

	for _, profaneWord := range profaneWords {
		loweredBody := strings.ToLower(cleanedBody)
		if strings.Contains(loweredBody, profaneWord) {
			words := strings.Split(cleanedBody, " ")
			for i, wordToSanitize := range words {
				if strings.ToLower(wordToSanitize) == profaneWord {
					words[i] = "****"
				}
			}
			cleanedBody = strings.Join(words, " ")
		}
	}

	chirp, err := api.DB.CreateChirp(cleanedBody)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}
