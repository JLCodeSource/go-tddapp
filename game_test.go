package poker_test

import (
	"testing"
	"github.com/vetch101/go-tddapp"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)


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

		assertFinishCalledWith(t, game, winner)
	})
	t.Run("start 3 player game, send blind alert on WS + finish with 'Chris' winner",
			 func(t *testing.T) {
		wantedBlindAlert := "Blind is 100"
		winner := "Chris"
		store := &poker.StubPlayerStore{}
		game := &GameSpy{BlindAlert: []byte(wantedBlindAlert)}
		
		playerServer := mustMakePlayerServer(t, store, game)

		server := httptest.NewServer(playerServer)
		ws := mustDialWS(t, "ws" + strings.TrimPrefix(server.URL, "http") + "/ws")

		defer server.Close()
		defer ws.Close()

		writeWSMessage(t, ws, "3")
		writeWSMessage(t, ws, winner)

		assertGameStartedWith(t, game, 3)
		assertFinishCalledWith(t, game, winner)

		timeout := (time.Duration(10) * time.Millisecond)
		within(t, timeout, func() {assertWebsocketGotMsg(t, ws, wantedBlindAlert)})
	})
}
