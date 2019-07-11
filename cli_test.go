package poker_test

import (
	"github.com/vetch101/go-tddapp"
	"strings"
	"testing"
	"bytes"
)

type GameSpy struct {
	StartedWith int
	FinishedWith string
}

func (g *GameSpy) Start(numberOfPlayers int) {
	g.StartedWith = numberOfPlayers
}

func (g *GameSpy) Finish(winner string) {
	g.FinishedWith = winner
}


func TestCLI(t *testing.T) {

	var dummyStdOut = &bytes.Buffer{}

	t.Run("it prompts the user to enter the number of players and starts the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		game := &GameSpy{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		gotPrompt := stdout.String()
		wantPrompt := poker.PlayerPrompt

		if gotPrompt != wantPrompt {
			t.Errorf("got '%s', want '%s'", gotPrompt, wantPrompt)
		}

		if game.StartedWith != 7 {
			t.Errorf("wanted Start called with 7 but got %d", game.StartedWith)
		}

	})

	t.Run("finish game with Chris as winner", func(t *testing.T) {

		in := strings.NewReader("1\nChris wins\n")
		game := &GameSpy{}

		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		if game.FinishedWith != "Chris" {
			t.Errorf("expected called with 'Chris' but got %q", game.FinishedWith)
		}

	})

	t.Run("record cleo win from user input", func(t *testing.T) {

		in := strings.NewReader("1\nCleo wins\n")
		game := &GameSpy{}
		
		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		if game.FinishedWith != "Cleo" {
			t.Errorf("expected called with 'Cleo' but got %q", game.FinishedWith)
		}
	})


}

