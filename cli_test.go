package poker_test

import (
	"github.com/vetch101/go-tddapp"
	"strings"
	"testing"
	"time"
	"fmt"
	"bytes"
)


// ScheduledAlert is a struct containing the time and amount of an alert
type ScheduledAlert struct {
	at time.Duration
	amount int
}

// String outputs the ScheduledAlert information
func (s ScheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

type SpyBlindAlerter struct {
	alerts []ScheduledAlert
}

func (s *SpyBlindAlerter) ScheduledAlertAt(duration time.Duration, amount int) {
	s.alerts = append(s.alerts, ScheduledAlert{duration, amount})
}

func TestCLI(t *testing.T) {

	
	var dummySpyAlerter = &SpyBlindAlerter{}
	var dummyPlayerStore = &poker.StubPlayerStore{}
	var dummyStdOut = &bytes.Buffer{}


	t.Run("record chris win from user input", func(t *testing.T) {

		in := strings.NewReader("5\nChris wins\n")
		playerStore := &poker.StubPlayerStore{}
		game := poker.NewGame(dummySpyAlerter, playerStore)

		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Chris")

	})

	t.Run("record cleo win from user input", func(t *testing.T) {

		in := strings.NewReader("5\nCleo wins\n")
		playerStore := &poker.StubPlayerStore{}
		game := poker.NewGame(dummySpyAlerter, playerStore)

		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Cleo")
	})

	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("5\nChris wins\n")
		playerStore := &poker.StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewGame(blindAlerter, playerStore)

		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

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

		for i, want := range cases {
			t.Run(fmt.Sprintf(want.String()), func(t *testing.T) {
					if len(blindAlerter.alerts) <= i {
						t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
					}

					got := blindAlerter.alerts[i]

					assertScheduledAlert(t, got, want)
				})
		}
	})

	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewGame(blindAlerter, dummyPlayerStore)

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		got := stdout.String()
		want := poker.PlayerPrompt

		if got != want {
			t.Errorf("got '%s', want '%s'", got, want)
		}

		cases := []ScheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]
				assertScheduledAlert(t, got, want)
			})
		}

	})
}

