package poker_test

import (
	"testing"
	"github.com/vetch101/go-tddapp"
	"os"
	"io/ioutil"
)


func TestFileSystemStore(t *testing.T) {

	database, cleanDatabase := createTempFile(t, "")
	defer cleanDatabase()

	t.Run("works with an empty file", func(t *testing.T) {

		_, err := poker.NewFileSystemPlayerStore(database)

		poker.AssertNoError(t, err)

	})

	database, cleanDatabase = createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
	defer cleanDatabase()

	store, err := poker.NewFileSystemPlayerStore(database)

	poker.AssertNoError(t, err)

	t.Run("/league from a reader sorted", func(t *testing.T) {

		got := store.GetLeague()

		want := []poker.Player{
			{Name: "Chris", Wins: 33},
			{Name: "Cleo", Wins: 10},
		}

		poker.AssertLeague(t, got, want)

		//read again
		got = store.GetLeague()
		poker.AssertLeague(t, got, want)
	})
	t.Run("get player score", func(t *testing.T) {

		got := store.GetPlayerScore("Chris")

		want := 33

		poker.AssertScoreEquals(t, got, want)
	})
	t.Run("store wins for existing players", func(t *testing.T) {
		store.PostRecordWin("Chris")

		got := store.GetPlayerScore("Chris")
		want := 34

		poker.AssertScoreEquals(t, got, want)
	})
	t.Run("store wins for new players", func(t *testing.T) {
		store.PostRecordWin("Joe")

		got := store.GetPlayerScore("Joe")
		want := 1
		poker.AssertScoreEquals(t, got, want)
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