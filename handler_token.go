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

	user, err := api.DB.GetUserByRefreshToken(authRefreshToken)

	if err != nil {
		log.Printf("Error with refresh token: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	if user.RefreshTokenExpiration.Compare(time.Now().UTC()) < 0 {
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

	err = api.DB.DeleteRefreshTokenByRefreshToken(authRefreshToken)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
