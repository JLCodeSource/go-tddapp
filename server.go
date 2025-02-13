package poker

import (
	"encoding/json"
	"fmt"
	websocket "github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const htmlTemplatePath = "game.html"

// Player stores a name with number of wins
type Player struct {
	Name string
	Wins int
}

// PlayerStore stores score information about players
type PlayerStore interface {
	GetPlayerScore(name string) int
	PostRecordWin(name string) error
	GetLeague() League
}

// PlayerServer is an HTTP interface for PlayerStore information
type PlayerServer struct {
	store PlayerStore
	http.Handler
	template *template.Template
	game     Game
}

type playerServerWS struct {
	*websocket.Conn
}

// NewPlayerServer instantiates a new PlayerServer
func NewPlayerServer(store PlayerStore, game Game) (*PlayerServer, error) {
	p := new(PlayerServer)

	tmpl, err := template.ParseFiles("game.html")

	if err != nil {
		return nil, fmt.Errorf("problem loading template %s %v", htmlTemplatePath, err)
	}

	p.template = tmpl
	p.store = store
	p.game = game

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))
	router.Handle("/game", http.HandlerFunc(p.gameHandler))
	router.Handle("/ws", http.HandlerFunc(p.webSocket))

	p.Handler = router

	return p, nil
}

func newPlayerServerWS(w http.ResponseWriter, r *http.Request) *playerServerWS {

	conn, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("problem upgrading connection to WebSocket %v\n", err)
	}

	return &playerServerWS{conn}

}

func (w *playerServerWS) WaitForMsg() string {
	_, msg, err := w.ReadMessage()
	if err != nil {
		log.Printf("error reading from websocket %v\n", err)
	}
	return string(msg)
}

func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {

	ws := newPlayerServerWS(w, r)

	numberOfPlayersMsg := ws.WaitForMsg()
	numberOfPlayers, _ := strconv.Atoi(string(numberOfPlayersMsg))
	p.game.Start(numberOfPlayers, ws)

	winnerMsg := ws.WaitForMsg()
	p.game.Finish(string(winnerMsg))

}

func (w *playerServerWS) Write(p []byte) (n int, err error) {
	err = w.WriteMessage(1, p)

	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func (p *PlayerServer) gameHandler(w http.ResponseWriter, r *http.Request) {
	p.template.Execute(w, nil)
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(p.store.GetLeague())
	w.WriteHeader(http.StatusOK)

}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/players/"):]

	switch r.Method {
	case http.MethodPost:
		p.postWin(w, player)
	case http.MethodGet:
		p.getScore(w, player)
	}
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
