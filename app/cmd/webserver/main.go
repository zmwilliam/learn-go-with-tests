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

	server := poker.NewPlayerServer(store)
	log.Fatal(http.ListenAndServe(":5000", server))
}
