package main

import (
	"log"
	"net/http"

	poker "github.com/zmwilliam/learn-go-with-tests/app"
)

const dbFileName = "game.db.json"

func main() {
	store, closeFn, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer closeFn()

	game := poker.NewTexasHoldem(poker.BlindAlerterFunc(poker.Alerter), store)

	server, err := poker.NewPlayerServer(store, game)
	if err != nil {
		log.Fatalf("problem creating player server %v", err)
	}

	log.Fatal(http.ListenAndServe(":5000", server))
}
