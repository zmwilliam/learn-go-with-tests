package poker_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zmwilliam/learn-go-with-tests/app"
)

func TestGETPlayers(t *testing.T) {
	store := poker.StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
	}
	server := newPlayerServer(t, &store)

	t.Run("return Pepper's score", func(t *testing.T) {
		request := poker.NewGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertResponseStatus(t, response.Code, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("return Floyd's score", func(t *testing.T) {
		request := poker.NewGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertResponseStatus(t, response.Code, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := poker.NewGetScoreRequest("missing_player")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertResponseStatus(t, response.Code, http.StatusNotFound)
	})
}

func newPlayerServer(t *testing.T, stubPlayerStore *poker.StubPlayerStore) *poker.PlayerServer {
	server, err := poker.NewPlayerServer(stubPlayerStore)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}

	return server
}

func TestStoreWins(t *testing.T) {
	store := poker.StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}

	server := newPlayerServer(t, &store)

	t.Run("it records wins when POST", func(t *testing.T) {
		player_name := "Pepper"
		request := poker.NewPostWinRequest(player_name)
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
		server := newPlayerServer(t, &store)

		request := poker.NewLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := poker.GetLeagueFromResponse(t, response.Body)
		poker.AssertResponseStatus(t, response.Code, http.StatusOK)
		poker.AssertLeague(t, got, wantedLeague)
		assertContentType(t, response, "application/json")
	})
}

func TestGame(t *testing.T) {
	t.Run("GET /game return 200", func(t *testing.T) {
		server := newPlayerServer(t, &poker.StubPlayerStore{})

		request := newGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		poker.AssertResponseStatus(t, response.Code, http.StatusOK)
	})

	t.Run("when we get a message over a websocket it is a winner of a game", func(t *testing.T) {
		store := &poker.StubPlayerStore{}
		winner := "Ruth"
		server := httptest.NewServer(newPlayerServer(t, store))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("could not open a ws connection on %s %v", wsURL, err)
		}
		defer ws.Close()

		if err := ws.WriteMessage(websocket.TextMessage, []byte(winner)); err != nil {
			t.Fatalf("could not message over ws connectino %v", err)
		}

		time.Sleep(10 * time.Millisecond)
		poker.AssertPlayerWin(t, store, winner)
	})

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
