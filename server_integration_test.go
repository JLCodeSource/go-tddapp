package poker_test

import (
	"github.com/vetch101/go-tddapp"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {

	database, cleanDatabase := createTempFile(t, `[]`)
	defer cleanDatabase()

	store, err := poker.NewFileSystemPlayerStore(database)
	poker.AssertNoError(t, err)

	server, _ := poker.NewPlayerServer(store, dummyGame)
	player := "Pepper"

	for i := 0; i < 3; i++ {
		server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	}

	t.Run("get score", func(t *testing.T) {

		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))

		poker.AssertStatus(t, response.Code, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), "3")
	})

	player = "Bob"

	for i := 0; i < 5; i++ {
		server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	}

	t.Run("get league", func(t *testing.T) {

		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())

		poker.AssertStatus(t, response.Code, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := []poker.Player{
			{Name: "Bob", Wins: 5},
			{Name: "Pepper", Wins: 3},
		}
		poker.AssertLeague(t, got, want)
	})

}
