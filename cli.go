package poker

import (
	"bufio"
	"io"
	"strings"
	"time"
)

// CLI is the playerstore and input reader for the commandline version
type CLI struct {
	playerStore PlayerStore
	in          *bufio.Scanner
	alerter		BlindAlerter
}

type BlindAlerter interface {
	ScheduledAlertAt(duration time.Duration, amount int)
}

// NewCLI is a constructor for playerStore
func NewCLI(store PlayerStore, in io.Reader, alerter BlindAlerter) *CLI {
	return &CLI{
		playerStore: store,
		in:          bufio.NewScanner(in),
		alerter:	alerter,
	}
}

// PlayPoker is the method to update the poker scores
func (cli *CLI) PlayPoker() {
	cli.alerter.ScheduledAlertAt(5 * time.Second, 100)
	userInput := cli.readLine()
	cli.playerStore.PostRecordWin(extractWinner(userInput))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
