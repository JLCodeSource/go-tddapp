package poker

import (
	"bufio"
	"io"
	"strings"
	"time"
	"fmt"
)

const PlayerPrompt = "Please enter the number of players: "

// CLI is the playerstore and input reader for the commandline version
type CLI struct {
	playerStore PlayerStore
	in          *bufio.Scanner
	out			io.Writer
	alerter		BlindAlerter
}

// NewCLI is a constructor for playerStore
func NewCLI(store PlayerStore, in io.Reader, out io.Writer, alerter BlindAlerter) *CLI {
	return &CLI{
		playerStore: store,
		in:          bufio.NewScanner(in),
		out: 		out,
		alerter:	alerter,
	}
}

// PlayPoker is the method to update the poker scores
func (cli *CLI) PlayPoker() {

	fmt.Fprint(cli.out, PlayerPrompt)

	cli.scheduleBlindAlerts()

	userInput := cli.readLine()
	cli.playerStore.PostRecordWin(extractWinner(userInput))
}

func (cli *CLI) scheduleBlindAlerts() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second

	for _, blind := range blinds {
		cli.alerter.ScheduledAlertAt(blindTime, blind)
		blindTime = blindTime + 10 * time.Minute	
	}
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
