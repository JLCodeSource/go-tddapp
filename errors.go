package poker

const (

	// ErrDBInitialize means that there was an error initilizing the db file
	ErrDBInitialize = Err("problem initilizing player db file")

	// ErrLoadingPlayerStore means there was an error loading the player store from file
	ErrLoadingPlayerStore = Err("problem loading player store from file")

	// ErrFileSeek means that there was an error seeking on file
	ErrFileSeek = Err("problem seeking on file")

	// ErrFileInfo means that there was an error getting file info from file
	ErrFileInfo = Err("problem getting file info from file")

	// ErrFileWrite means that there was an error writing to file
	ErrFileWrite = Err("problem writing to file")

	// ErrFileOpen means that there was an error opening the file
	ErrFileOpen = Err("problem opening file")

	// ErrCreateStore means there was an error creating the player store
	ErrCreateStore = Err("problem creating file system player store")

	// ErrFileClose means there was an error closing the file
	ErrFileClose = Err("problem closing file")

	// ErrEncode means there was an error during json encoding
	ErrEncode = Err("problem encoding json")

	// ErrBadPlayerInput is an error for bad inputs 
	ErrBadPlayerInput = "Bad value received for number of players, please try again with a number"
)

// Err are errors that can happen when interacting with FSPS
type Err string

// The Error func returns the FSPS Error
func (e Err) Error() string {
	return string(e)
}
