package poker

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {

	database, cleanDatabase := createTempFile(t, `[]`)
	defer cleanDatabase()
	
	store, err := NewFileSystemPlayerStore(database)
	assertNoError(t, err)

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
			{"Bob", 5},
			{"Pepper", 3},
		}
		assertLeague(t, got, want)
	})
	
}