package poker

import (
	"bufio"
	"io"
	"strings"
	"fmt"
	"strconv"
)

// CLI is the playerstore and input reader for the commandline version
type CLI struct {
	//PlayerStore PlayerStore
	in          *bufio.Scanner
	out			io.Writer
	game		Game
}

// PlayerPrompt is the prompt for number of players
const PlayerPrompt = "Please enter the number of players: "

// BadWinnerInputMsg is the prompt for a bad winner input
const BadWinnerInputMsg = "You entered an incorrect value. Please enter '{Playername} wins'"

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
	numberOfPlayers, err := strconv.Atoi(strings.Trim(numberOfPlayersInput, "\n"))

	if err != nil {
		fmt.Fprint(cli.out, ErrBadPlayerInput)
		return
	}

	cli.game.Start(numberOfPlayers, cli.out)

	winnerInput := cli.readLine()

	if strings.Contains(winnerInput, " wins") == false {
		fmt.Fprint(cli.out, BadWinnerInputMsg)
	}

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
