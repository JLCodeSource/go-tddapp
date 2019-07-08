package poker

import (
	"testing"
	"net/http/httptest"
	"reflect"
)

const jsonContentType = "application/json"

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league League
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) GetLeague() League {
	return s.league
}

func (s *StubPlayerStore) PostRecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func AssertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("response status is wrong - got status %d, want %d", got, want)
	}
}

func AssertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong - got '%s', want '%s'", got, want)
	}
}

func AssertLeague(t *testing.T, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func AssertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	
	t.Helper()

	got := response.Result().Header.Get("content-type")

	if got != want {
		t.Errorf("response did not have content-type of application/json, got %v, want %v",
			got, want)
	}

} 

func AssertPlayerWin(t *testing.T, store *StubPlayerStore, winner string) {
	t.Helper()

	gotLen := len(store.winCalls)
	wantLen := 1

	if gotLen != wantLen {
		t.Fatalf("got %d calls to RecordWin want %d", gotLen, wantLen)
	}

	gotWinner := store.winCalls[0]
	if gotWinner != winner {
		t.Errorf("did not store correct winner got '%s' want '%s", 
			gotWinner, winner)
	}
}

func AssertScoreEquals(t *testing.T, got, want int) {
	
	t.Helper()

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func AssertNoError(t *testing.T, err error) {

	if err != nil {
		t.Fatalf("didn't expect error but got one, %v", err)
	}

}