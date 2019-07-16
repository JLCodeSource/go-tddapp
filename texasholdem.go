package poker

import (
	"time"
	"os"
	"io"
)

// TexasHoldEm is a struct containing alerter (a BlindAlerter) 
// and store (PlayerStore)
type TexasHoldEm struct {
	alerter BlindAlerter
	store PlayerStore
}

// NewTexasHoldEm returns a pointer to a TexasHoldEm struct
func NewTexasHoldEm(alerter BlindAlerter, store PlayerStore, to io.Writer) *TexasHoldEm {
	return &TexasHoldEm{
		alerter:alerter,
		store:store,
	}
}

// Start starts a game of TexasHoldEm with numberOfPlayers
func (t *TexasHoldEm) Start(numberOfPlayers int) {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	blindIncrement := time.Duration(5 + numberOfPlayers) * time.Minute

	for _, blind := range blinds {
		t.alerter.ScheduledAlertAt(blindTime, blind, os.Stdout)
		blindTime = blindTime + blindIncrement
	}
}

// Finish finishes the game of TexasHoldEm recording the winner
func (t *TexasHoldEm) Finish(winner string) {
	t.store.PostRecordWin(winner)
}