package poker_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	poker "github.com/zmwilliam/learn-go-with-tests/app"
)

var dummyGame = &GameSpy{}

func TestGETPlayers(t *testing.T) {
	store := poker.StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
	}
	server := mustCreatePlayerServer(t, &store, dummyGame)

	t.Run("return Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertResponseStatus(t, response.Code, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("return Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertResponseStatus(t, response.Code, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("missing_player")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertResponseStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestStoreWins(t *testing.T) {
	store := poker.StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}

	server := mustCreatePlayerServer(t, &store, dummyGame)

	t.Run("it records wins when POST", func(t *testing.T) {
		player_name := "Pepper"
		request := newPostWinRequest(player_name)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertResponseStatus(t, response.Code, http.StatusAccepted)

		got := len(store.WinCalls)
		want := 1

		if got != want {
			t.Errorf("got %d calls to RecordWin want %d", got, want)
		}

		got_winner := store.WinCalls[0]

		if got_winner != player_name {
			t.Errorf("did not store correct winner got %q want %q", got_winner, player_name)
		}
	})
}

func TestLeague(t *testing.T) {
	t.Run("it return the league table as JSON", func(t *testing.T) {
		wantedLeague := poker.League{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		store := poker.StubPlayerStore{nil, nil, wantedLeague}
		server := mustCreatePlayerServer(t, &store, dummyGame)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)
		poker.AssertResponseStatus(t, response.Code, http.StatusOK)
		poker.AssertLeague(t, got, wantedLeague)
		assertContentType(t, response, "application/json")
	})
}

func TestGame(t *testing.T) {
	t.Run("GET /game return 200", func(t *testing.T) {
		server := mustCreatePlayerServer(t, &poker.StubPlayerStore{}, dummyGame)

		request := newGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		poker.AssertResponseStatus(t, response.Code, http.StatusOK)
	})

	t.Run("start a game with 3 players, send some blind alerts down WS and declare Ruth the winner", func(t *testing.T) {
		wantedBlindAlert := "Blind is 100"
		winner := "Ruth"

		game := &GameSpy{BlindAlert: []byte(wantedBlindAlert)}
		dummyStore := &poker.StubPlayerStore{}
		server := httptest.NewServer(mustCreatePlayerServer(t, dummyStore, game))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
		ws := mustDialWS(t, wsURL)
		defer ws.Close()

		writeWSMessage(t, ws, "3")
		writeWSMessage(t, ws, winner)

		assertGameStartedWith(t, game, 3)
		assertFinishCalledWith(t, game, winner)

		within(t, 10*time.Millisecond, func() {
			assertWebsocketGotMsg(t, ws, wantedBlindAlert)
		})
	})

}

func assertWebsocketGotMsg(t *testing.T, ws *websocket.Conn, want string) {
	_, gotMsg, _ := ws.ReadMessage()

	if string(gotMsg) != want {
		t.Errorf("got %s, want %s", string(gotMsg), want)
	}
}

func writeWSMessage(t *testing.T, ws *websocket.Conn, message string) {
	t.Helper()

	err := ws.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		t.Fatalf("could not message over ws connectino %v", err)
	}
}

func mustDialWS(t *testing.T, wsURL string) *websocket.Conn {
	t.Helper()

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", wsURL, err)
	}
	return ws
}

func newGameRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/game", nil)
	return req
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	got_content_type := response.Result().Header.Get("content-type")
	if got_content_type != want {
		t.Errorf("response did not have content-type of %s, got %v", want, got_content_type)
	}
}

func mustCreatePlayerServer(t *testing.T, store poker.PlayerStore, game poker.Game) *poker.PlayerServer {
	server, err := poker.NewPlayerServer(store, game)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}

	return server
}

func newPostWinRequest(name string) *http.Request {
	url_path := fmt.Sprintf("/players/%s", name)
	req, _ := http.NewRequest(http.MethodPost, url_path, nil)
	return req
}

func newGetScoreRequest(name string) *http.Request {
	url_path := fmt.Sprintf("/players/%s", name)
	req, _ := http.NewRequest(http.MethodGet, url_path, nil)
	return req
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func getLeagueFromResponse(t testing.TB, body io.Reader) (league poker.League) {
	t.Helper()

	err := json.NewDecoder(body).Decode(&league)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v", body, err)
	}

	return
}

func within(t testing.TB, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)

	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Error("timed out")
	case <-done:
	}

}
