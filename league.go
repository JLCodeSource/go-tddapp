package main

import (
	"fmt"
	"io"
	"encoding/json"
)

// League is an array of Players
type League []Player

// NewLeague is a constructor for a League
func NewLeague(rdr io.Reader) ([]Player, error) {
	
	var l League
	err := json.NewDecoder(rdr).Decode(&l)

	if err != nil {
		err = fmt.Errorf("problem parsing league, %v", err)
	}

	return l, err
}

// Find finds and returns a Player
func (l League) Find(name string) *Player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}
	return nil
}