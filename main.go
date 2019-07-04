package main

import (
	"log"
	"net/http"
)

// InMemoryPlayerStore is a struct for holding player scores
type InMemoryPlayerStore struct{}

// GetPlayerScore queries a name on an InMemPlayStr and returns score
func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return 123
}

// RecordWin increments a win
func (i *InMemoryPlayerStore) RecordWin(name string) {}

func main() {
	server := &PlayerServer{&InMemoryPlayerStore{}}
	
	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}