package poker

import (
	"encoding/json"
	"os"
	"fmt"
	"sort"
)

// FileSystemPlayerStore is a json.Encoder that stores a League of Players[]
type FileSystemPlayerStore struct {
	database *json.Encoder
	league League
}

func initializePlayerDBFile(file *os.File) error {

	_, err := file.Seek(0,0)
	
	if err != nil {
		return fmt.Errorf("problem seeking on file %s, %v",
			file.Name(), err)
	}

	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v",
			file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0,0)
	}

	return nil

}

// NewFileSystemPlayerStore is a constructor method for the FielSystemPlayerStore
func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
	
	err := initializePlayerDBFile(file)

	if err != nil {
		return nil, fmt.Errorf("problem initilizing player db file, %v", err)
	}

	league, err := NewLeague(file)

	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", 
			file.Name(), err)
	}
	
	return &FileSystemPlayerStore{
		database:json.NewEncoder(&tape{file}),
		league:league,
	}, nil
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
func (f *FileSystemPlayerStore) PostRecordWin(name string) {

	player := f.league.Find(name)

	if player != nil {
			player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1})
	}

	f.database.Encode(f.league)

}