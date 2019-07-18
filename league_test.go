package poker_test

import (
	"github.com/vetch101/go-tddapp"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLeague(t *testing.T) {

	t.Run("it returns league table as JSON", func(t *testing.T) {

		wantedLeague := []poker.Player{
			{Name: "Cleo", Wins: 32},
			{Name: "Chris", Wins: 20},
			{Name: "Trevor", Wins: 12},
		}

		store := poker.StubPlayerStore{Scores: nil, WinCalls: nil, League: wantedLeague}
		server, _ := poker.NewPlayerServer(&store, dummyGame)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		league := getLeagueFromResponse(t, response.Body)

		contentType := response.Result().Header.Get("content-type")

		poker.AssertStatus(t, response.Code, http.StatusOK)

		poker.AssertContentType(t, contentType, jsonContentType)

		poker.AssertLeague(t, league, wantedLeague)
	})
}
