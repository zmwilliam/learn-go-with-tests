package poker

import (
	"fmt"
	"io"
	"reflect"
	"testing"
	"time"
)

type StubPlayerStore struct {
	Scores   map[string]int
	WinCalls []string
	League   League
}

func (s *StubPlayerStore) GetLeague() League {
	return s.League
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.Scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.WinCalls = append(s.WinCalls, name)
}

type SpyBlindAlerter struct {
	Alerts []ScheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int, to io.Writer) {
	s.Alerts = append(s.Alerts, ScheduledAlert{at, amount})
}

type ScheduledAlert struct {
	At     time.Duration
	Amount int
}

func (s ScheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.Amount, s.At)
}

func AssertPlayerWin(t testing.TB, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.WinCalls) != 1 {
		t.Fatalf("expected a win call but did't get any")
	}

	got := store.WinCalls[0]
	if got != winner {
		t.Errorf("didn't record correct winner, got %q, want %q", got, winner)
	}

}

func AssertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}

func AssertResponseStatus(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got status %d want %d", got, want)
	}
}

func AssertResponseBody(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func AssertLeague(t testing.TB, got, want League) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
