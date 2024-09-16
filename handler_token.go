package main

import (
	"log"
	"net/http"
	"time"

	"github.com/iamhectorsosa/web-server/internal/auth"
)

func (api *apiConfig) postRefresh(w http.ResponseWriter, r *http.Request) {
	authRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthenticated request")
		return
	}

	user, refreshToken, err := api.DB.GetUserAndRefreshTokenByRefreshToken(authRefreshToken)

	if err != nil {
		log.Printf("Error with refresh token: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired")
		return
	}

	token, err := auth.CreateJWT(user.Id, api.jwtSecret, 0)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "JWT Token Creation failed")
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}

func (api *apiConfig) postRevoke(w http.ResponseWriter, r *http.Request) {

	authRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthenticated request")
		return
	}

	err = api.DB.DeleteRefreshToken(authRefreshToken)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid refresh token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
