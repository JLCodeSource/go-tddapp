package poker_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"strings"
	websocket "github.com/gorilla/websocket"
	"time"
	"github.com/vetch101/go-tddapp"
)

// jsonContentType refers to the JSON http content header
const jsonContentType = "application/json"

var (
	dummyGame = &GameSpy{}
)
func TestGETPlayers(t *testing.T) {

	store := poker.StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		}, nil, nil,
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
		map[string]int{},
		nil, nil,
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

func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server, _ := poker.NewPlayerServer(&poker.StubPlayerStore{}, dummyGame)

		request := newGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response.Code, http.StatusOK)
	})
	t.Run("when we get a message over a websocket it is a winner", func(t *testing.T) {
		store := &poker.StubPlayerStore{}
		winner := "Ruth"
		game := dummyGame
		playerServer := mustMakePlayerServer(t, store, dummyGame)
		server := httptest.NewServer(playerServer)
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
		ws := mustDialWS(t, wsURL)
		defer ws.Close()

		writeWSMessage(t, ws, "3")
		writeWSMessage(t, ws, winner)

		time.Sleep(100 * time.Millisecond)
		assertFinishCalledWith(t, game, winner)
	})
	t.Run("start game with 3 players and finish game with 'Chris' as winner", func(t *testing.T) {
		store := &poker.StubPlayerStore{}
		game := dummyGame
		winner := "Chris"
		playerServer := mustMakePlayerServer(t, store, game)

		server := httptest.NewServer(playerServer)
		ws := mustDialWS(t, "ws" + strings.TrimPrefix(server.URL, "http") + "/ws")

		defer server.Close()
		defer ws.Close()

		writeWSMessage(t, ws, "3")
		writeWSMessage(t, ws, winner)

		time.Sleep(10* time.Millisecond)
		assertGameStartedWith(t, game, 3)
		assertFinishCalledWith(t, game, winner)
	})
}

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

func TestFileSystemStore(t *testing.T) {

	database, cleanDatabase := createTempFile(t, "")
	defer cleanDatabase()

	t.Run("works with an empty file", func(t *testing.T) {

		_, err := poker.NewFileSystemPlayerStore(database)

		poker.AssertNoError(t, err)

	})

	database, cleanDatabase = createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
	defer cleanDatabase()

	store, err := poker.NewFileSystemPlayerStore(database)

	poker.AssertNoError(t, err)

	t.Run("/league from a reader sorted", func(t *testing.T) {

		got := store.GetLeague()

		want := []poker.Player{
			{"Chris", 33},
			{"Cleo", 10},
		}

		poker.AssertLeague(t, got, want)

		//read again
		got = store.GetLeague()
		poker.AssertLeague(t, got, want)
	})
	t.Run("get player score", func(t *testing.T) {

		got := store.GetPlayerScore("Chris")

		want := 33

		poker.AssertScoreEquals(t, got, want)
	})
	t.Run("store wins for existing players", func(t *testing.T) {
		store.PostRecordWin("Chris")

		got := store.GetPlayerScore("Chris")
		want := 34

		poker.AssertScoreEquals(t, got, want)
	})
	t.Run("store wins for new players", func(t *testing.T) {
		store.PostRecordWin("Joe")

		got := store.GetPlayerScore("Joe")
		want := 1
		poker.AssertScoreEquals(t, got, want)
	})

}

func createTempFile(t *testing.T, initialData string) (*os.File, func()) {
	t.Helper()

	tmpFile, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpFile.Write([]byte(initialData))

	removeFile := func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}

	return tmpFile, removeFile
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
	if game.StartCalled == false {
		t.Errorf("game should have started but did not start")
	}
	got := game.StartedWith
	if got != want {
		t.Errorf("got %d players, but wanted %d", got, want)
	}
}

func assertFinishCalledWith(t *testing.T, game *GameSpy, want string) {
	t.Helper()
	if game.FinishCalled == false {
		t.Errorf("game should have finished but did not finish")
	}
	got := game.FinishedWith 
	if got != want {
		t.Errorf("got %s winner, but wanted %s", got, want)
	}
}