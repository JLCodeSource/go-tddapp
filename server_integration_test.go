package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {

	database, cleanDatabase := createTempFile(t, "")
	defer cleanDatabase()
	store := &FileSystemPlayerStore{database}
	server := NewPlayerServer(store)
	player := "Pepper"

	for i:=0; i<3; i++ {
		server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	}


	t.Run("get score", func(t *testing.T) {

		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		
		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "3")
	})

	player = "Bob"

	for i:=0; i<5; i++ {
		server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	}

	t.Run("get league", func(t *testing.T) {

		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())

		assertStatus(t, response.Code, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := []Player{
			{"Pepper", 3},
			{"Bob", 5},
		}
		assertLeague(t, got, want)
	})
	
}