package poker

import (
	"time"
)

type TexasHoldEm struct {
	alerter BlindAlerter
	store PlayerStore
}

func NewTexasHoldEm(alerter BlindAlerter, store PlayerStore) *TexasHoldEm {
	return &TexasHoldEm{
		alerter:alerter,
		store:store,
	}
}

func (t *TexasHoldEm) Start(numberOfPlayers int) {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	blindIncrement := time.Duration(5 + numberOfPlayers) * time.Minute

	for _, blind := range blinds {
		t.alerter.ScheduledAlertAt(blindTime, blind)
		blindTime = blindTime + blindIncrement
	}
}

func (t *TexasHoldEm) Finish(winner string) {
	t.store.PostRecordWin(winner)
}