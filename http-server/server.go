package main

import (
	"fmt"
	"net/http"
	"strings"
)

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
}

type PlayerServer struct {
	store PlayerStore
}

func (s *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player_name := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodGet:
		s.getScore(w, player_name)
	case http.MethodPost:
		s.postScore(w, player_name)
	}
}

func (s *PlayerServer) getScore(w http.ResponseWriter, player_name string) {
	score := s.store.GetPlayerScore(player_name)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (s *PlayerServer) postScore(w http.ResponseWriter, player_name string) {
	s.store.RecordWin(player_name)
	w.WriteHeader(http.StatusAccepted)
}
