package poker

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
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
	store    PlayerStore
	template *template.Template
	game     Game
	http.Handler
}

const htmlTemplatePath = "game.html"

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewPlayerServer(store PlayerStore, game Game) (*PlayerServer, error) {
	server := new(PlayerServer)

	tmpl, err := template.ParseFiles("game.html")
	if err != nil {
		return nil, fmt.Errorf("problem opening %s %v", htmlTemplatePath, err)
	}

	server.template = tmpl
	server.store = store
	server.game = game

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(server.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(server.playersHandler))
	router.Handle("/game", http.HandlerFunc(server.gameHandler))
	router.Handle("/ws", http.HandlerFunc(server.webSocketHandler))

	server.Handler = router

	return server, nil
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

func (s *PlayerServer) gameHandler(w http.ResponseWriter, r *http.Request) {
	s.template.Execute(w, nil)
}

func (s *PlayerServer) webSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, _ := wsUpgrader.Upgrade(w, r, nil)

	_, numberOfPlayersMsg, _ := conn.ReadMessage()
	numberOfPlayers, _ := strconv.Atoi(string(numberOfPlayersMsg))
	s.game.Start(numberOfPlayers, io.Discard) //TODO we still discanting blinds messages

	_, winnerMsg, _ := conn.ReadMessage()
	s.game.Finish(string(winnerMsg))
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
