package main

import (
	"fmt"
	"blackjack/deck"
)

func main() {
	cardDeck := deck.NewDefaultDeck()
	cardDeck.Shuffle(5)
	fmt.Println(cardDeck.Inspect())
}
