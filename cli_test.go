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

}

