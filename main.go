package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	database "github.com/iamhectorsosa/web-server/database"
)

const port = "8080"
const filepathRoot = "."

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	mux := http.NewServeMux()

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/*", fsHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("GET /api/reset", apiCfg.handleReset)

	mux.HandleFunc("GET /api/healthz", handleReadiness)

	// mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.getChirpById)
	mux.HandleFunc("POST /api/chirps", apiCfg.postChirps)

	mux.HandleFunc("POST /api/users", apiCfg.postUsers)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: http://localhost:%s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DB.GetChirps()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) getChirpById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	chirpId, err := strconv.Atoi(id)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID")
		return
	}

	chirps, err := cfg.DB.GetChirpById(chirpId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) postChirps(w http.ResponseWriter, r *http.Request) {
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

	chirp, err := cfg.DB.CreateChirp(cleanedBody)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) postUsers(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Email string `json:"email"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	user, err := cfg.DB.CreateUser(payload.Email)

	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits++
		w.Header().Add("Control-cache", "no-cache")
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
<html>

<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>

</html>
	`, cfg.fileserverHits)))
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}
