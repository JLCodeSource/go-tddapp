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

	return nil

}

// PostRecordWin increments a win
func (i *InMemoryPlayerStore) PostRecordWin(name string) {
	i.store[name]++
}
