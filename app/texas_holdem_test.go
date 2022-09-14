package poker_test

import (
	"fmt"
	"testing"
	"time"

	poker "github.com/zmwilliam/learn-go-with-tests/app"
)

func TestGame_Start(t *testing.T) {
	var dummyPlayerStore = &poker.StubPlayerStore{}

	t.Run("schedules alerts on game start for 5 players", func(t *testing.T) {
		blindAlerter := &poker.SpyBlindAlerter{}

		game := poker.NewTexasHoldem(blindAlerter, dummyPlayerStore)
		game.Start(5)

		cases := []poker.ScheduledAlert{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		checkSchedulingCases(t, cases, blindAlerter)
	})

	t.Run("schedules alerts on game start for a 7 players", func(t *testing.T) {
		blindAlerter := &poker.SpyBlindAlerter{}

		game := poker.NewTexasHoldem(blindAlerter, dummyPlayerStore)
		game.Start(7)

		cases := []poker.ScheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		checkSchedulingCases(t, cases, blindAlerter)
	})
}

func TestGame_Finish(t *testing.T) {
	dummyBlindAlerter := &poker.SpyBlindAlerter{}
	store := &poker.StubPlayerStore{}
	winner := "Ruth"

	game := poker.NewTexasHoldem(dummyBlindAlerter, store)
	game.Finish(winner)

	poker.AssertPlayerWin(t, store, winner)
}

func checkSchedulingCases(t *testing.T, cases []poker.ScheduledAlert, blindAlerter *poker.SpyBlindAlerter) {
	t.Helper()

	for i, want := range cases {
		t.Run(fmt.Sprint(want), func(t *testing.T) {
			if len(blindAlerter.Alerts) <= 1 {
				t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.Alerts)
			}

			got := blindAlerter.Alerts[i]
			assertScheduledAlert(t, got, want)
		})
	}
}

func assertScheduledAlert(t *testing.T, got, want poker.ScheduledAlert) {
	t.Helper()
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
