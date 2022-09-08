package main

import (
	"io"
)

type FileSystemPlayerStore struct {
	database io.ReadSeeker
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	league := f.GetLeague()
	for _, player := range league {
		if player.Name == name {
			return player.Wins
		}
	}

	return 0
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	//TODO not implemented yet
}

func (f *FileSystemPlayerStore) GetLeague() []Player {
	f.database.Seek(0, 0)
	league, _ := NewLeague(f.database)
	return league
}
