package poker_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/zmwilliam/learn-go-with-tests/app"
)

type GameSpy struct {
	StartedWith  int
	FinishedWith string
}

func (g *GameSpy) Start(numberOfPlayers int) {
	g.StartedWith = numberOfPlayers
}

func (g *GameSpy) Finish(winner string) {
	g.FinishedWith = winner
}

func TestCLI(t *testing.T) {
	var dummyStdOut = &bytes.Buffer{}

	t.Run("it prompts the user to enter the number of players and starts the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		game := &GameSpy{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		got := stdout.String()
		want := poker.PlayerPrompt

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		if game.StartedWith != 7 {
			t.Errorf("want started with 7 but got %d", game.StartedWith)
		}
	})

	t.Run("Finish game with correct winner", func(t *testing.T) {
		for _, winner := range []string{"Chris", "Cleo"} {
			in := strings.NewReader(fmt.Sprintf("1\n%s wins\n", winner))
			game := &GameSpy{}
			cli := poker.NewCLI(in, dummyStdOut, game)
			cli.PlayPoker()

			if game.FinishedWith != winner {
				t.Errorf("expected finish called with %s but got %q", winner, game.FinishedWith)
			}
		}
	})

	t.Run("Do not read beyond the first newline", func(t *testing.T) {
		dummyAlerter := &poker.SpyBlindAlerter{}
		dummyStore := &poker.StubPlayerStore{}
		in := failOnEndReader{
			t,
			strings.NewReader("1\nChris win\nhello there"),
		}

		game := poker.NewTexasHoldem(dummyAlerter, dummyStore)
		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()
	})
}

type failOnEndReader struct {
	t      *testing.T
	reader io.Reader
}

func (m failOnEndReader) Read(p []byte) (n int, err error) {
	n, err = m.reader.Read(p)

	if n == 0 || err == io.EOF {
		m.t.Fatal("Read to the end when you should not have")
	}

	return n, err
}
