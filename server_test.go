package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

func TestGETPlayers(t *testing.T) {
	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		PlayerServer(response, request)

		got := response.Body.String()
		want := "20"

		assertResponseBody(t, got, want)
		
	})
	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		PlayerServer(response, request)

		got := response.Body.String()
		want := "10"

		assertResponseBody(t, got, want)
	})
}

func newGetScoreRequest(name string) *http.Request {
	path := "/players/" + name
	req, _ := http.NewRequest(http.MethodGet, path, nil)
	return req
}

func assertResponseBody(t *testing.T, got, want string){
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}