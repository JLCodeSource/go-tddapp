package main

// NewInMemoryPlayerStore initializes an empty player store
func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}}
}

// InMemoryPlayerStore is a struct for holding player scores
type InMemoryPlayerStore struct{
	store map[string]int
}

// GetPlayerScore queries a name and returns score
func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.store[name]
}

// GetLeague
func (i *InMemoryPlayerStore) GetLeague() []Player {

	var league []Player
	for name, wins := range i.store {
		league = append(league, Player{name, wins})
	}
	return league

}

// PostRecordWin increments a win
func (i *InMemoryPlayerStore) PostRecordWin(name string) {
	i.store[name]++
}
