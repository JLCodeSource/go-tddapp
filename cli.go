package poker

import (
	"io"
	"strings"
	"bufio"
)

// CLI is the playerstore and input reader for the commandline version
type CLI struct {
	playerStore PlayerStore
	in *bufio.Scanner
}

// NewCLI is a constructor for playerStore
func NewCLI(store PlayerStore, in io.Reader) *CLI {
	return &CLI {
		playerStore: store,
		in: bufio.NewScanner(in),
	}
}

// PlayPoker is the method to update the poker scores
func (cli *CLI) PlayPoker() {
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