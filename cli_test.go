package poker_test

import (
	"github.com/vetch101/go-tddapp"
	"strings"
	"testing"
	"bytes"
	"io"
)

type GameSpy struct {
	StartCalled bool
	StartedWith int
	BlindAlert []byte

	FinishCalled bool
	FinishedWith string
}

func (g *GameSpy) Start(numberOfPlayers int, out io.Writer) {
	g.StartCalled = true
	g.StartedWith = numberOfPlayers
	out.Write(g.BlindAlert)
}

func (g *GameSpy) Finish(winner string) {
	g.FinishCalled = true
	g.FinishedWith = winner
}

func userSends(messages ...string) io.Reader {
	return strings.NewReader(strings.Join(messages, "\n"))
}

func TestCLI(t *testing.T) {

	var dummyStdOut = &bytes.Buffer{}

	t.Run("it prompts the user to enter the number of players and starts the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := userSends("7", "Bob wins")
		game := &GameSpy{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		wantPrompt := poker.PlayerPrompt

		assertMessageSentToUser(t, stdout, wantPrompt)
		assertGameStarted(t, game.StartCalled)
		assertNumberOfPlayers(t, game.StartedWith, 7)

	})

	t.Run("finish game with Chris as winner", func(t *testing.T) {

		in := userSends("1", "Chris wins")
		game := &GameSpy{}

		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		assertGameFinished(t, game.FinishCalled)
		assertGameWonBy(t, "Chris", game.FinishedWith)
		
	})

	t.Run("prints error on non-numeric value entered + does not start", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := userSends("blah")
		game := &GameSpy{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		wantPrompt := poker.PlayerPrompt + poker.ErrBadPlayerInput

		assertGameNotStarted(t, game.StartCalled)
		assertMessageSentToUser(t, stdout, wantPrompt)
	})

	t.Run("prints an error if non-name wins entered", func (t *testing.T) {
		stdout := &bytes.Buffer{}
		in := userSends("1", "Lloyd is a killer")
		game := &GameSpy{}
		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()
		
		assertGameStarted(t, game.StartCalled)
		assertMessageSentToUser(t, stdout, poker.PlayerPrompt, poker.BadWinnerInputMsg)

	})
}

func assertMessageSentToUser(t *testing.T, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()
	if got != want {
		t.Errorf("got '%s' sent to stdout but expected '%+v'", got, want)
	}
}

func assertGameStarted(t *testing.T, started bool) {
	t.Helper()
	if started != true {
		t.Errorf("game should have started but did not start")
	}
}

func assertGameNotStarted(t *testing.T, started bool) {
	t.Helper()
	if started == true {
		t.Errorf("game should not have started but has started")
	}
}

func assertNumberOfPlayers(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("wanted Start called with %d but got %d", got, want)
	}
}

func assertGameFinished(t *testing.T, finished bool) {
	t.Helper()
	if finished != true {
		t.Errorf("game should have finished but did not finish")
	}
}

func assertGameWonBy(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("expected called with '%s' but got '%s'", got, want)
	}
}