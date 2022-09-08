package main

import (
	"os"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	database, cleanDatabaseFn := createTempFile(t, `[
		  {"Name": "Cleo", "Wins": 10},
		  {"Name": "Chris", "Wins": 33}]`)

	defer cleanDatabaseFn()

	store := NewFileSystemStore(database)

	t.Run("league from a reader", func(t *testing.T) {
		got := store.GetLeague()
		want := []Player{
			{"Cleo", 10},
			{"Chris", 33},
		}

		assertLeague(t, got, want)
	})

	t.Run("get player score", func(t *testing.T) {
		got := store.GetPlayerScore("Chris")
		want := 33
		assertScoreEquals(t, got, want)
	})

	t.Run("store wins for existing players", func(t *testing.T) {
		winner_name := "Chris"
		store.RecordWin(winner_name)
		got := store.GetPlayerScore(winner_name)
		want := 34
		assertScoreEquals(t, got, want)
	})

	t.Run("store wins for new players", func(t *testing.T) {
		winner_name := "Pepper"
		store.RecordWin(winner_name)

		got := store.GetPlayerScore(winner_name)
		want := 1

		assertScoreEquals(t, got, want)
	})
}

func assertScoreEquals(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "db")
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFIleFn := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFIleFn
}
