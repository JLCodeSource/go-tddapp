package poker

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"reflect"
	"io"
	"io/ioutil"
	"os"
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

func TestGETPlayers(t *testing.T) {

	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		}, nil, nil,
	}
	server := NewPlayerServer(&store)

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "20")

	})
	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "10")
	})
	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestPOSTWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		nil, nil,
	}
	server := NewPlayerServer(&store)

	t.Run("it records win on POST", func(t *testing.T) {
		player := "Pepper"

		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to PostRecordWin want %d", len(store.winCalls), 1)
		}

		if store.winCalls[0] != player {
			t.Errorf("did not store correct winner - got '%s' want '%s'",
				store.winCalls[0], player)
		}
	})

}

func TestLeague(t *testing.T) {
	
	t.Run("it returns league table as JSON", func(t *testing.T) {
	
		wantedLeague := []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Trevor", 12},
		}
	
		store := StubPlayerStore{nil, nil, wantedLeague}
		server := NewPlayerServer(&store)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		assertContentType(t, response, jsonContentType)
		
		got := getLeagueFromResponse(t, response.Body)
				
		assertLeague(t, got, wantedLeague)
	})
}

func TestFileSystemStore(t *testing.T) {
	
	database, cleanDatabase := createTempFile(t, "")
	defer cleanDatabase()

	t.Run("works with an empty file", func(t *testing.T) {

		_, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)

	})

	database, cleanDatabase = createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
	defer cleanDatabase()

	store, err := NewFileSystemPlayerStore(database)
	
	assertNoError(t, err)

	t.Run("/league from a reader sorted", func(t *testing.T) {
		t.Helper()

		got := store.GetLeague()

		want := []Player{
			{"Chris", 33},
			{"Cleo", 10},
		}

		assertLeague(t, got, want)

		//read again
		got = store.GetLeague()
		assertLeague(t, got, want)
	})
	t.Run("get player score", func(t *testing.T) {
		t.Helper()

		got := store.GetPlayerScore("Chris")

		want := 33

		assertScoreEquals(t, got, want)
	})
	t.Run("store wins for existing players", func(t *testing.T) {
		store.PostRecordWin("Chris")

		got := store.GetPlayerScore("Chris")
		want := 34

		assertScoreEquals(t, got, want)
	})
	t.Run("store wins for new players", func(t *testing.T) {
		store.PostRecordWin("Joe")

		got := store.GetPlayerScore("Joe")
		want := 1
		assertScoreEquals(t, got, want)
	})
	
}

func createTempFile(t *testing.T, initialData string) (*os.File, func()) {
	t.Helper()

	tmpFile, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpFile.Write([]byte(initialData))

	removeFile := func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}

	return tmpFile, removeFile
}

func newGetScoreRequest(name string) *http.Request {
	path := "/players/" + name
	req, _ := http.NewRequest(http.MethodGet, path, nil)
	return req
}

func newPostWinRequest(name string) *http.Request {
	path := "/players/" + name
	req, _ := http.NewRequest(http.MethodPost, path, nil)
	return req
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func getLeagueFromResponse(t *testing.T, body io.Reader) []Player {
	t.Helper()
	league, _ := NewLeague(body)
	return league
} 

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("response status is wrong - got status %d, want %d", got, want)
	}
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong - got '%s', want '%s'", got, want)
	}
}

func assertLeague(t *testing.T, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	
	t.Helper()

	got := response.Result().Header.Get("content-type")

	if got != want {
		t.Errorf("response did not have content-type of application/json, got %v, want %v",
			got, want)
	}

} 

func assertScoreEquals(t *testing.T, got, want int) {
	
	t.Helper()

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func assertNoError(t *testing.T, err error) {

	if err != nil {
		t.Fatalf("didn't expect error but got one, %v", err)
	}

}