package poker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Player struct {
	Name string
	Wins int
}

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() League
}

type PlayerServer struct {
	store PlayerStore
	http.Handler
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	server := new(PlayerServer)
	server.store = store

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(server.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(server.playersHandler))

	server.Handler = router

	return server
}

func (s *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(s.store.GetLeague())
}

func (s *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
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
