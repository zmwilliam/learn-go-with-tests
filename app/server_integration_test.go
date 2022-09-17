package poker_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	poker "github.com/zmwilliam/learn-go-with-tests/app"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := createTempFile(t, `[]`)
	defer cleanDatabase()
	store, err := poker.NewFileSystemStore(database)

	poker.AssertNoError(t, err)

	server := mustCreatePlayerServer(t, store, dummyGame)

	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		poker.AssertResponseStatus(t, response.Code, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		poker.AssertResponseStatus(t, response.Code, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := []poker.Player{{"Pepper", 3}}

		poker.AssertLeague(t, got, want)
	})

}
