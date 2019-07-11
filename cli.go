package poker

import (
	"bufio"
	"io"
	"strings"
	"time"
	"fmt"
	"strconv"
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

	numberOfPlayers, _ := strconv.Atoi(cli.readLine())

	cli.scheduleBlindAlerts(numberOfPlayers)

	userInput := cli.readLine()
	cli.playerStore.PostRecordWin(extractWinner(userInput))
}

func (cli *CLI) scheduleBlindAlerts(numberOfPlayers int) {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	blindIncrement := time.Duration(5 + numberOfPlayers) * time.Minute

	for _, blind := range blinds {
		cli.alerter.ScheduledAlertAt(blindTime, blind)
		blindTime = blindTime + blindIncrement
	}
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
