package main

import (
	"fmt"
	"io"
	"encoding/json"
)

type League []Player

func NewLeague(rdr io.Reader) ([]Player, error) {
	
	var l League
	err := json.NewDecoder(rdr).Decode(&l)

	if err != nil {
		err = fmt.Errorf("problem parsing league, %v", err)
	}

	return l, err
}

func (l League) Find(name string) *Player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}
	return nil
}