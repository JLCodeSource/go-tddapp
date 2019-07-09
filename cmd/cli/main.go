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
	poker.NewCLI(store, os.Stdin).PlayPoker()
}