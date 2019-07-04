package main

import (
	"net/http"
	"fmt"
)

// PlayerStore is an interface implementing GetPlyrScore
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
}

// PlayerServer is a struct with a store representing PlayerStore
type PlayerServer struct {
	store PlayerStore
}


// PlayerServer is a http server that provides the http.Method routing
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
	switch r.Method {
	case http.MethodPost:
		p.processWin(w, r)
	case http.MethodGet:
		p.showScore(w, r)
	}
}

func (p *PlayerServer) showScore(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/players/"):]

	score := p.store.GetPlayerScore(player)

	if score == 0 {
	w.WriteHeader(http.StatusNotFound)
	}
	
	fmt.Fprint(w, score)

}

func (p *PlayerServer) processWin(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/players/"):]

	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}

