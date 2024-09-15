package main

import (
	"flag"
	"log"
	"os"

	"github.com/iamhectorsosa/web-server/internal/database"
	"github.com/joho/godotenv"
)

const port = "8080"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	dbg := flag.Bool("debug", false, "Enable debug mode and get a fresh database to start with.")
	flag.Parse()

	databaseStore, err := database.NewDB("database.json", *dbg)
	if err != nil {
		log.Fatal(err)
	}

	api := apiConfig{DB: databaseStore, jwtSecret: jwtSecret}
	server := NewServer(api, port)
	log.Printf("Listening on port: http://localhost:%s\n", port)
	log.Fatal(server.ListenAndServe())
}
