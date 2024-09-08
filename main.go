package main

import (
	"flag"
	"log"

	"github.com/iamhectorsosa/web-server/database"
)

const port = "8080"

func main() {
	dbg := flag.Bool("debug", false, "Enable debug mode and get a fresh database to start with.")
	flag.Parse()

	databaseStore, err := database.NewDB("database.json", *dbg)
	if err != nil {
		log.Fatal(err)
	}

	api := apiConfig{DB: databaseStore}
	server := NewServer(api, port)
	log.Printf("Listening on port: http://localhost:%s\n", port)
	log.Fatal(server.ListenAndServe())
}
