package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/iamhectorsosa/web-server/internal/auth"
	database "github.com/iamhectorsosa/web-server/internal/database"
)

func (api *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	authorIdStr := r.URL.Query().Get("author_id")

	if authorIdStr == "" {
		authorIdStr = "0"
	}

	authorId, err := strconv.Atoi(authorIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Author ID")
		return
	}

	sortQ := r.URL.Query().Get("sort")

	if sortQ != "asc" && sortQ != "desc" {
		sortQ = ""
	}

	dbChirps, err := api.DB.GetChirps()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve Chirps")
		return
	}

	chirps := []database.Chirp{}
	for _, dbChirp := range dbChirps {
		if authorId != 0 && dbChirp.AuthorId != authorId {
			continue
		}

		chirps = append(chirps, database.Chirp{
			Id:       dbChirp.Id,
			AuthorId: dbChirp.AuthorId,
			Body:     dbChirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		if sortQ == "desc" {
			return chirps[i].Id > chirps[j].Id
		}
		return chirps[i].Id < chirps[j].Id
	})

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

	decoder := json.NewDecoder(r.Body)
	payload := struct {
		Body string `json:"body"`
	}{}

	err = decoder.Decode(&payload)

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

	chirp, err := api.DB.CreateChirp(cleanedBody, userId)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func (api *apiConfig) deleteChirpById(w http.ResponseWriter, r *http.Request) {
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

	id := r.PathValue("id")
	chirpId, err := strconv.Atoi(id)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID")
		return
	}

	chirp, err := api.DB.GetChirpById(chirpId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Something went wrong")
		return
	}

	if chirp.AuthorId != userId {
		respondWithError(w, http.StatusForbidden, "Cannot delete others Chirps")
		return
	}

	err = api.DB.DeleteChirpById(chirpId)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
