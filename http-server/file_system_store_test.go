package main

import (
	"strings"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	database := strings.NewReader(`[
		  {"Name": "Cleo", "Wins": 10},
		  {"Name": "Chris", "Wins": 33}]`)

	store := FileSystemPlayerStore{database}

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

		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})
}
