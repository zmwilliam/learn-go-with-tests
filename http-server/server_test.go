package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func TestGETPlayers(t *testing.T) {
	store := &StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
	}
	server := &PlayerServer{store: store}

	t.Run("return Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("return Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("missing_player")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, http.StatusNotFound)

	})
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		nil,
	}

	server := &PlayerServer{&store}

	// t.Run("it returns accepted on POST", func(t *testing.T) {
	// 	request := newPostWinRequest("Pepper")
	// 	response := httptest.NewRecorder()
	//
	// 	server.ServeHTTP(response, request)
	//
	// 	assertResponseStatus(t, response.Code, http.StatusAccepted)
	// })

	t.Run("it records wins when POST", func(t *testing.T) {
		player_name := "Pepper"
		request := newPostWinRequest(player_name)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, http.StatusAccepted)

		got := len(store.winCalls)
		want := 1

		if got != want {
			t.Errorf("got %d calls to RecordWin want %d", got, want)
		}

		got_winner := store.winCalls[0]

		if got_winner != player_name {
			t.Errorf("did not store correct winner got %q want %q", got_winner, player_name)
		}
	})
}

func newGetScoreRequest(name string) *http.Request {
	url_path := fmt.Sprintf("/players/%s", name)
	req, _ := http.NewRequest(http.MethodGet, url_path, nil)
	return req
}

func newPostWinRequest(name string) *http.Request {
	url_path := fmt.Sprintf("/players/%s", name)
	req, _ := http.NewRequest(http.MethodPost, url_path, nil)
	return req
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertResponseStatus(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got status %d want %d", got, want)
	}
}
