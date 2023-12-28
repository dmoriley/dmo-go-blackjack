package main

import (
	// "blackjack/deck"
	// "blackjack/engine"
	"blackjack/game"
	"fmt"
	"os"
	"os/user"
)

func main() {
	// cardDeck := engine.NewBlackjackDeck(2)
	// cardDeck.Shuffle(5)
	// fmt.Println(cardDeck.Inspect())

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("The name of the game is Blackjack\n")
	fmt.Printf("Don't be caught counting cards %s...\n", user.Username)
	game.Start(os.Stdin, os.Stdout)
}
