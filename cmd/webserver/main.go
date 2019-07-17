package main

import (
	"github.com/vetch101/go-tddapp"
	"log"
	"net/http"
)

const dbFileName = "game.db.json"

func main() {

	store, close, err := poker.FileSystemStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	alerter := poker.BlindAlerterFunc(poker.Alerter)
	game := poker.NewTexasHoldEm(alerter, store)

	server, _ := poker.NewPlayerServer(store, game)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
