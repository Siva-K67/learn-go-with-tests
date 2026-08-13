package main

import (
	"log"
	"net/http"
	"os"
)

const dbFileName = "game.db.json"

// main opens (or creates) the database file, wires it into the server, and starts listening
func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666) // open/create the JSON db file, read+write
	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := NewFileSystemPlayerStore(db) // wraps the file as our PlayerStore
	if err != nil {
		log.Fatalf("problem initializing player store %v", err)
	}
	server := NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
