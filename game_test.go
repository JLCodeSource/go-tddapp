package poker_test

import (
	"github.com/vetch101/go-tddapp"
	"testing"
	"time"
	"fmt"
)

var dummySpyAlerter = &SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}

func TestGame_Start(t *testing.T) {
	t.Run("it schedules blind values for 5 players", func(t *testing.T) {
		
		blindAlerter := &SpyBlindAlerter{}
		
		game := poker.NewGame(blindAlerter, dummyPlayerStore)

		game.Start(5)
		
		cases := []ScheduledAlert {
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		checkSchedulingCases(cases, t, blindAlerter)
		
	})

	t.Run("it schedules alerts for 7 players", func(t *testing.T) {
		
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewGame(blindAlerter, dummyPlayerStore)

		game.Start(7)
		
		cases := []ScheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		checkSchedulingCases(cases, t, blindAlerter)

	})
}

func Test_Finish(t *testing.T) {
	store := &poker.StubPlayerStore{}
	game := poker.NewGame(dummySpyAlerter, store)
	winner := "Ruth"

	game.Finish(winner)
	poker.AssertPlayerWin(t, store, winner)
}

func checkSchedulingCases(cases []ScheduledAlert, t *testing.T, alerter poker.BlindAlerter) {
	t.Helper()
	for i, want := range cases {
		t.Run(fmt.Sprintf(want.String()), func(t *testing.T) {
				if len(cases) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, cases)
				}

				got := cases[i]

				assertScheduledAlert(t, got, want)
			})
	}
}

func assertScheduledAlert(t *testing.T, got, want ScheduledAlert) {
	t.Helper()

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}