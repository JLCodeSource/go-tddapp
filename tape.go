package poker

import "os"

// Tape is a struct for an os.File
type Tape struct {
	File *os.File
}

// Write is a writer for the tape file
func (t *Tape) Write(p []byte) (n int, err error) {
	t.File.Truncate(0)
	t.File.Seek(0, 0)
	return t.File.Write(p)
}
