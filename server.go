package main

import (
	"net/http"

	"github.com/iamhectorsosa/web-server/database"
)

type apiConfig struct {
	DB *database.DB
}

func NewServer(api apiConfig, port string) *http.Server {
	router := http.NewServeMux()
	router.HandleFunc("GET /api/chirps", api.getChirps)
	router.HandleFunc("GET /api/chirps/{id}", api.getChirpById)
	router.HandleFunc("POST /api/chirps", api.postChirps)

	router.HandleFunc("POST /api/users", api.postUsers)
	router.HandleFunc("POST /api/login", api.postLogin)

	return &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
}
