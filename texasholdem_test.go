package poker_test

import (
	"fmt"
	"github.com/vetch101/go-tddapp"
	"io"
	"os"
	"testing"
	"time"
)

var dummySpyAlerter = &SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}

// ScheduledAlert is a struct containing the time and amount of an alert
type ScheduledAlert struct {
	at     time.Duration
	amount int
}

// String outputs the ScheduledAlert information
func (s ScheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

type SpyBlindAlerter struct {
	alerts []ScheduledAlert
}

func (s *SpyBlindAlerter) ScheduledAlertAt(duration time.Duration, amount int, to io.Writer) {
	s.alerts = append(s.alerts, ScheduledAlert{duration, amount})
}

func TestGame_Start(t *testing.T) {
	t.Run("it schedules blind values for 5 players", func(t *testing.T) {

		blindAlerter := &SpyBlindAlerter{}
		alertsDestination := os.Stdout

		game := poker.NewTexasHoldEm(blindAlerter, dummyPlayerStore)

		game.Start(5, alertsDestination)

		cases := []ScheduledAlert{
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
		alertsDestination := os.Stdout
		game := poker.NewTexasHoldEm(blindAlerter, dummyPlayerStore)

		game.Start(7, alertsDestination)

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
	alertsDestination := os.Stdout
	game := poker.NewTexasHoldEm(dummySpyAlerter, store)
	winner := "Ruth"
	game.Start(1, alertsDestination)
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
