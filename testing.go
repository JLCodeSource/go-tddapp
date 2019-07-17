package poker

import (
	"reflect"
	"testing"
)

// StubPlayerStore is a spy stub mock for PlayerStore
type StubPlayerStore struct {
	Scores   map[string]int
	WinCalls []string
	League   []Player
}

// GetPlayerScore returns the spy store score
func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.Scores[name]
	return score
}

// GetLeague returns the spy store league of Players[]
func (s *StubPlayerStore) GetLeague() League {
	return s.League
}

// PostRecordWin adds to the wins in winCalls
func (s *StubPlayerStore) PostRecordWin(name string) error {
	s.WinCalls = append(s.WinCalls, name)
	return nil
}

// AssertStatus is an assertion for http response status
func AssertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("response status is wrong - got status %d, want %d", got, want)
	}
}

// AssertResponseBody is an assertion for http response body
func AssertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong - got '%s', want '%s'", got, want)
	}
}

// AssertContentType is an assertion for the http ressponse content-type
func AssertContentType(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response did not have content-type of application/json, got %v, want %v",
			got, want)
	}
}

// AssertLeague asserts the content of the League
func AssertLeague(t *testing.T, got, want League) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

// AssertPlayerWin asserts which player won (& that it only wins once)
func AssertPlayerWin(t *testing.T, store *StubPlayerStore, winner string) {
	t.Helper()

	gotLen := len(store.WinCalls)
	wantLen := 1

	if gotLen != wantLen {
		t.Fatalf("got %d calls to RecordWin want %d", gotLen, wantLen)
	}

	gotWinner := store.WinCalls[0]
	if gotWinner != winner {
		t.Errorf("did not store correct winner got '%s' want '%s",
			gotWinner, winner)
	}
}

// AssertScoreEquals asserts the score
func AssertScoreEquals(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

// AssertNoError asserts that there is no error
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect error but got one, %v", err)
	}
}