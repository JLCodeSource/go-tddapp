package poker

import (
	"time"
)

type Game struct {
	alerter BlindAlerter
	store PlayerStore
}

func NewGame(alerter BlindAlerter, store PlayerStore) *Game {
	return &Game{
		alerter:alerter,
		store:store,
	}
}

func (g *Game) Start(numberOfPlayers int) {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	blindIncrement := time.Duration(5 + numberOfPlayers) * time.Minute

	for _, blind := range blinds {
		g.alerter.ScheduledAlertAt(blindTime, blind)
		blindTime = blindTime + blindIncrement
	}
}

func (g *Game) Finish(winner string) {
	g.store.PostRecordWin(winner)
}