package main

import (
	"net/http"
	"fmt"
)

// PlayerStore is an interface implementing GetPlyrScore
type PlayerStore interface {
	GetPlayerScore(name string) int
	//PostPlayerScore(name string) int
}

// PlayerServer is a struct with a store representing PlayerStore
type PlayerServer struct {
	store PlayerStore
}


// PlayerServer is a http server that returns the player score
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/players/"):]
	
	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusAccepted)
		return
	}
	
	score := p.store.GetPlayerScore(player)

	if score == 0 {
	w.WriteHeader(http.StatusNotFound)
	}
	
	fmt.Fprint(w, score)
	
}

