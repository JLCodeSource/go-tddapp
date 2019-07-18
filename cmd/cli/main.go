package main

import (
	"fmt"
	"github.com/vetch101/go-tddapp"
	"log"
	"os"
)

const dbFileName = "game.db.json"

func main() {

	store, close, err := poker.FileSystemStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")
	alerter := poker.BlindAlerterFunc(poker.Alerter)

	game := poker.NewTexasHoldEm(alerter, store)
	cli := poker.NewCLI(os.Stdin, os.Stdout, game)
	cli.PlayPoker()
}
