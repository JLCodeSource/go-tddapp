package main

import (
	"fmt"
	"net/http"
)

// PlayerStore is an interface implementing GetPlyrScore
type PlayerStore interface {
	GetPlayerScore(name string) int
	PostRecordWin(name string)
}

// PlayerServer is a struct with a store representing PlayerStore
type PlayerServer struct {
	store PlayerStore
}

// PlayerServer is a http server that provides the http.Method routing
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	router := http.NewServeMux()

	router.Handle("/league", 
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

	router.Handle("/players/", 
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			player := r.URL.Path[len("/players/"):]

			switch r.Method {
			case http.MethodPost:
				p.postWin(w, player)
			case http.MethodGet:
				p.getScore(w, player)
			}
		}))
	
	router.ServeHTTP(w, r)
}

func (p *PlayerServer) getScore(w http.ResponseWriter, player string) {

	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)

}

func (p *PlayerServer) postWin(w http.ResponseWriter, player string) {

	p.store.PostRecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}
