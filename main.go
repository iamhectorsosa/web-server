package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const port = "8080"
const filepathRoot = "."

type apiConfig struct {
	fileserverHits int
}

func main() {
	mux := http.NewServeMux()
	apiCfg := apiConfig{}

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/*", fsHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("GET /api/reset", apiCfg.handleReset)

	mux.HandleFunc("GET /api/healthz", handleReadiness)

	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)

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

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	// decoder.DisallowUnknownFields()
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

	sanitizedChirp := payload.Body

	for _, word := range profaneWords {
		lowerCaseChirp := strings.ToLower(sanitizedChirp)
		hasProfane := strings.Contains(lowerCaseChirp, word)
		if hasProfane {
			arrayChirp := strings.Split(sanitizedChirp, " ")
			var sanitizedArrayChirp []string
			for _, wordToSanitize := range arrayChirp {
				if strings.ToLower(wordToSanitize) == word {
					sanitizedArrayChirp = append(sanitizedArrayChirp, "****")
				} else {
					sanitizedArrayChirp = append(sanitizedArrayChirp, wordToSanitize)
				}
			}
			sanitizedChirp = strings.Join(sanitizedArrayChirp, " ")
		}
	}

	respondWithJSON(w, http.StatusOK, struct {
		CleanedBody string `json:"cleaned_body"`
	}{
		CleanedBody: sanitizedChirp,
	})

}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	res, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(res)
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
