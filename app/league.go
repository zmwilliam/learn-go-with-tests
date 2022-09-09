package poker

import (
	"encoding/json"
	"fmt"
	"io"
)

type League []Player

func (league League) Find(name string) *Player {
	for i, p := range league {
		if p.Name == name {
			return &league[i]
		}
	}

	return nil
}

func NewLeague(reader io.Reader) (League, error) {
	var league League
	err := json.NewDecoder(reader).Decode(&league)
	if err != nil {
		err = fmt.Errorf("problem parsing league, %v", err)
	}
	return league, err
}
