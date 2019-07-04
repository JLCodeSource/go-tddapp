package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := PlayerServer{store}
	player := "Pepper"

	for i:=0; i<3; i++ {
		server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	}

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(player))
	
	assertStatus(t, response.Code, http.StatusOK)
	assertResponseBody(t, response.Body.String(), "3")
}