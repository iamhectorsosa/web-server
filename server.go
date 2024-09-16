package main

import (
	"net/http"

	"github.com/iamhectorsosa/web-server/internal/database"
)

type apiConfig struct {
	DB        *database.DB
	jwtSecret string
}

func NewServer(api apiConfig, port string) *http.Server {
	router := http.NewServeMux()
	router.HandleFunc("GET /api/chirps", api.getChirps)
	router.HandleFunc("GET /api/chirps/{id}", api.getChirpById)
	router.HandleFunc("POST /api/chirps", api.postChirps)
	router.HandleFunc("DELETE /api/chirps/{id}", api.deleteChirpById)

	router.HandleFunc("POST /api/users", api.postUsers)
	router.HandleFunc("PUT /api/users", api.putUsers)
	router.HandleFunc("POST /api/login", api.postLogin)

	router.HandleFunc("POST /api/refresh", api.postRefresh)
	router.HandleFunc("POST /api/revoke", api.postRevoke)

	return &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
}
