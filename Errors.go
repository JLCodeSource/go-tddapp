package poker

const (

	// ErrDBInitialize means that there was an error initilizing the db file
	ErrDBInitialize = PokerErr("problem initilizing player db file")

	// ErrLoadingPlayerStore means there was an error loading the player store from file
	ErrLoadingPlayerStore = PokerErr("problem loading player store from file")

	// ErrFileSeek means that there was an error seeking on file
	ErrFileSeek = PokerErr("problem seeking on file")

	// ErrFileInfo means that there was an error getting file info from file
	ErrFileInfo = PokerErr("problem getting file info from file")

	// ErrFileWrite means that there was an error writing to file
	ErrFileWrite = PokerErr("problem writing to file")

	// ErrFileOpen means that there was an error opening the file
	ErrFileOpen = PokerErr("problem opening file")

	// ErrCreateStore means there was an error creating the player store
	ErrCreateStore = PokerErr("problem creating file system player store")

	// ErrFileClose means there was an error closing the file
	ErrFileClose = PokerErr("problem closing file")
)

// PokerErr are errors that can happen when interacting with FSPS
type PokerErr string

// The Error func returns the FSPS Error
func (e PokerErr) Error() string {
	return string(e)
}
