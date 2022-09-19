package poker_test

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/zmwilliam/learn-go-with-tests/app"
)

type GameSpy struct {
	StartCalled bool
	StartedWith int
	BlindAlert  []byte

	FinishCalled bool
	FinishedWith string
}

func (g *GameSpy) Start(numberOfPlayers int, out io.Writer) {
	g.StartCalled = true
	g.StartedWith = numberOfPlayers
	out.Write(g.BlindAlert)
}

func (g *GameSpy) Finish(winner string) {
	g.FinishCalled = true
	g.FinishedWith = winner
}

func TestCLI(t *testing.T) {
	var dummyStdOut = &bytes.Buffer{}

	t.Run("start game with 3 players and finish game with 'Chris' as winner", func(t *testing.T) {
		game := &GameSpy{}
		stdout := &bytes.Buffer{}

		in := userSends("3", "Chris wins")
		cli := poker.NewCLI(in, stdout, game)

		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt)
		assertGameStartedWith(t, game, 3)
		assertFinishCalledWith(t, game, "Chris")
	})

	t.Run("start game with 8 players and record 'Cleo' as winner", func(t *testing.T) {
		game := &GameSpy{}

		in := userSends("8", "Cleo wins")
		cli := poker.NewCLI(in, dummyStdOut, game)

		cli.PlayPoker()

		assertGameStartedWith(t, game, 8)
		assertFinishCalledWith(t, game, "Cleo")
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

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		stout := &bytes.Buffer{}
		in := strings.NewReader("Pies\n")
		game := &GameSpy{}

		cli := poker.NewCLI(in, stout, game)
		cli.PlayPoker()

		if game.StartCalled {
			t.Errorf("game should not have started")
		}

		assertMessagesSentToUser(t, stout, poker.PlayerPrompt, poker.BadPlayerInputErrMsg)
	})

	t.Run("it prints an error when winner is declared incorrectly", func(t *testing.T) {
		stout := &bytes.Buffer{}
		in := strings.NewReader("2\nLloyd is a Killer\n")
		game := &GameSpy{}

		cli := poker.NewCLI(in, stout, game)
		cli.PlayPoker()

		if game.FinishCalled {
			t.Errorf("game should not have finished")
		}

		assertMessagesSentToUser(t, stout, poker.PlayerPrompt, poker.BadWinnerInputErrMsg)
	})
}

func assertFinishCalledWith(t *testing.T, game *GameSpy, wantedWinner string) {
	t.Helper()

	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.FinishedWith == wantedWinner
	})

	if !passed {
		t.Errorf("expected finish called with %q but got %q", wantedWinner, game.FinishedWith)
	}
}

func assertGameStartedWith(t *testing.T, game *GameSpy, numberOfPlayersWanted int) {
	t.Helper()

	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.StartedWith == numberOfPlayersWanted
	})

	if !passed {
		t.Errorf("expcted Start called with %d but got %d ", numberOfPlayersWanted, game.StartedWith)
	}
}

func assertMessagesSentToUser(t testing.TB, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()
	if got != want {
		t.Errorf("got %q sent to stdout but expected %+v", got, messages)
	}
}

func userSends(messages ...string) io.Reader {
	return strings.NewReader(strings.Join(messages, "\n"))
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

func retryUntil(d time.Duration, f func() bool) bool {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if f() {
			return true
		}
	}
	return false
}
