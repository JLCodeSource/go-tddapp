package main

import (
	"log"
	"net/http"
	"github.com/vetch101/go-tddapp"
)

const dbFileName = "game.db.json"

func main() {

	store, close, err := poker.FileSystemStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()
	
	server := poker.NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
