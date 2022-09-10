package main

import (
	"fmt"
	"log"
	"os"

	poker "github.com/zmwilliam/learn-go-with-tests/app"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	store, closeFn, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer closeFn()

	game := poker.NewCLI(store, os.Stdin)
	game.PlayPoker()
}