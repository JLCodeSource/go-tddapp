package poker

import (
	"bufio"
	"io"
	"strings"
	"fmt"
	"strconv"
)

// Game interface is what starts and finishes games
type Game interface {
	Start(numberOfPlayers int)
	Finish(winner string)
}

// CLI is the playerstore and input reader for the commandline version
type CLI struct {
	PlayerStore PlayerStore
	in          *bufio.Scanner
	out			io.Writer
	game		Game
}

// PlayerPrompt is the prompt for number of players
const PlayerPrompt = "Please enter the number of players: "

// NewCLI is a constructor for playerStore
func NewCLI(in io.Reader, out io.Writer, game Game) *CLI {
	return &CLI{
		in:          bufio.NewScanner(in),
		out: 		out,
		game: game,
	}
}

// PlayPoker is the method to update the poker scores
func (cli *CLI) PlayPoker() {

	fmt.Fprint(cli.out, PlayerPrompt)

	numberOfPlayersInput := cli.readLine()
	numberOfPlayers, _ := strconv.Atoi(strings.Trim(numberOfPlayersInput, "\n"))

	cli.game.Start(numberOfPlayers)

	winnerInput := cli.readLine()
	winner := extractWinner(winnerInput)

	cli.game.Finish(winner)
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
