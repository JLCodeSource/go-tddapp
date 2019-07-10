package poker

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// FileSystemPlayerStore is a json.Encoder that stores a League of Players[]
type FileSystemPlayerStore struct {
	database *json.Encoder
	league   League
}

// NewFileSystemPlayerStore is a constructor method for the FielSystemPlayerStore
func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {

	err := initializePlayerDBFile(file)

	if err != nil {
		return nil, ErrDBInitialize
	}

	league, err := NewLeague(file)

	if err != nil {
		return nil, ErrLoadingPlayerStore
	}

	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}, nil
}

func initializePlayerDBFile(file *os.File) error {

	_, err := file.Seek(0, 0)

	if err != nil {
		return ErrFileSeek
	}

	info, err := file.Stat()

	if err != nil {
		return ErrFileInfo
	}

	if info.Size() == 0 {
		_, err := file.Write([]byte("[]"))
		if err != nil {
			return ErrFileWrite
		}
		_, err = file.Seek(0, 0)
		if err != nil {
			return ErrFileSeek
		}
	}

	return nil

}

// FileSystemStoreFromFile retrieves a file and returns the FSPS, its close func and the error
func FileSystemStoreFromFile(filename string) (*FileSystemPlayerStore, func(), error) {

	db, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0600)

	if err != nil {
		return nil, nil, ErrFileOpen
	}

	closeFunc := func() {
		e := db.Close()
		if e != nil {
			fmt.Errorf(string(ErrFileClose))
		}
	}
	store, err := NewFileSystemPlayerStore(db)

	if err != nil {
		return nil, nil, ErrCreateStore
	}

	return store, closeFunc, nil

}

// GetLeague is a method on a FSPlayerStore that sorts the League
func (f *FileSystemPlayerStore) GetLeague() League {

	sort.Slice(f.league, func(i, j int) bool {
		return f.league[i].Wins > f.league[j].Wins
	})
	return f.league
}

// GetPlayerScore returns a player's score (or zero if they don't exist)
func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {

	player := f.GetLeague().Find(name)

	if player != nil {
		return player.Wins
	}

	return 0
}

// PostRecordWin increments a player's score (or creates the player if they don't exist)
func (f *FileSystemPlayerStore) PostRecordWin(name string) error {

	player := f.league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1})
	}

	err := f.database.Encode(f.league)

	if err != nil {
		return ErrEncode
	}

	return nil
}
