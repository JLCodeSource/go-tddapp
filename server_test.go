package poker_test

import (
	websocket "github.com/gorilla/websocket"
	"github.com/vetch101/go-tddapp"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// jsonContentType refers to the JSON http content header
const jsonContentType = "application/json"

var (
	dummyGame = &GameSpy{}
)

func TestGETPlayers(t *testing.T) {

	store := poker.StubPlayerStore{
		Scores: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		WinCalls: nil,
		League:   nil,
	}
	server, _ := poker.NewPlayerServer(&store, dummyGame)

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response.Code, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), "20")

	})
	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response.Code, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), "10")
	})
	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestPOSTWins(t *testing.T) {
	store := poker.StubPlayerStore{
		Scores:   map[string]int{},
		WinCalls: nil, League: nil,
	}
	server, _ := poker.NewPlayerServer(&store, dummyGame)

	t.Run("it records win on POST", func(t *testing.T) {
		player := "Pepper"

		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response.Code, http.StatusAccepted)

		poker.AssertPlayerWin(t, &store, "Pepper")
	})

}

func mustDialWS(t *testing.T, url string) *websocket.Conn {
	t.Helper()
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", url, err)
	}
	return ws
}

func mustMakePlayerServer(t *testing.T, store poker.PlayerStore, game poker.Game) *poker.PlayerServer {
	t.Helper()
	server, err := poker.NewPlayerServer(store, game)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}

	return server
}

func writeWSMessage(t *testing.T, conn *websocket.Conn, message string) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}

func newGetScoreRequest(name string) *http.Request {
	path := "/players/" + name
	req, _ := http.NewRequest(http.MethodGet, path, nil)
	return req
}

func newPostWinRequest(name string) *http.Request {
	path := "/players/" + name
	req, _ := http.NewRequest(http.MethodPost, path, nil)
	return req
}

func newGameRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/game", nil)
	return request
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func getLeagueFromResponse(t *testing.T, body io.Reader) []poker.Player {
	t.Helper()
	league, _ := poker.NewLeague(body)
	return league
}

func assertGameStartedWith(t *testing.T, game *GameSpy, want int) {
	t.Helper()

	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.StartedWith == want
	})

	if game.StartCalled == false {
		t.Errorf("game should have started but did not start")
	}
	got := game.StartedWith
	if !passed {
		t.Errorf("got %d players, but wanted %d", got, want)
	}
}

func assertFinishCalledWith(t *testing.T, game *GameSpy, winner string) {
	t.Helper()

	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.FinishedWith == winner
	})

	if !passed {
		t.Errorf("expected finish called, but finish was not called")
	}

	if game.FinishCalled == false {
		t.Errorf("game should have finished but did not finish")
	}

	got := game.FinishedWith
	if got != winner {
		t.Errorf("got %s winner, but wanted %s", got, winner)
	}

}

func assertWebsocketGotMsg(t *testing.T, ws *websocket.Conn, want string) {
	_, got, _ := ws.ReadMessage()

	if string(got) != want {
		t.Errorf("got blind alert %q, want %q", string(got), string(want))
	}
}

func within(t *testing.T, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)

	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Error("timed out")
	case <-done:
	}
}

func retryUntil(d time.Duration, f func() bool) bool {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if f() {
			return true
		}
	}
	return false
}
