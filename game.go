package poker

import (
	"io"
)

// Game interface is what starts and finishes games within the CLI
type Game interface {
	Start(numberOfPlayers int, alertsDestination io.Writer)
	Finish(winner string)
}
