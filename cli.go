package poker

import (
	"io"
	"strings"
	"bufio"
)

// CLI is the playerstore and input reader for the commandline version
type CLI struct {
	playerStore PlayerStore
	in io.Reader
}

// PlayPoker is the method to update the poker scores
func (cli *CLI) PlayPoker() {
	reader := bufio.NewScanner(cli.in)
	reader.Scan()
	cli.playerStore.PostRecordWin(extractWinner(reader.Text()))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}