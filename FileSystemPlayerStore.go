package main

import (
	"io"
	"encoding/json"
)

type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
}

func (f *FileSystemPlayerStore) GetLeague() League {
	f.database.Seek(0,0)
	league, _ := NewLeague(f.database)
	return league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	
	player := f.GetLeague().Find(name)
	
	if player != nil {
		return player.Wins
	}

	return 0
}

func (f *FileSystemPlayerStore) PostRecordWin(name string) {

	league := f.GetLeague()
	player := league.Find(name)

	if player != nil {
			player.Wins++
	}

	f.database.Seek(0,0)
	json.NewEncoder(f.database).Encode(league)

}