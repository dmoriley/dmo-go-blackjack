package game

import (
	"blackjack/card/rank"
	"blackjack/deck"
	"bytes"
	"fmt"
)

type BlackjackDeck struct {
	DeckCount int
	deck.Deck
}

func NewBlackjackDeck(numberOfDecks int) *BlackjackDeck {
	bjDeck := &BlackjackDeck{
		DeckCount: 0,
	}

	deckToAdd := deck.NewDeck()

	// change all the kings, queens and jacks to value of 10
	for _, card := range deckToAdd.Cards {
		if card.Rank.Name == rank.King || card.Rank.Name == rank.Queen ||
			card.Rank.Name == rank.Jack {
			card.Rank.Value = 10
		}
	}

	for i := 0; i < numberOfDecks; i++ {
		bjDeck.AddCards(deckToAdd.Cards)
		bjDeck.DeckCount++
	}

	return bjDeck
}

func (d *BlackjackDeck) Inspect() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("*** Deck of %d cards ***\n", d.GetLength()))
	out.WriteString(fmt.Sprintf("Deck count: %d\n", d.DeckCount))
	out.WriteString(fmt.Sprintf("Shuffle count: %d\n", d.ShuffleCount))

	return out.String()
}

func (d *BlackjackDeck) Debug() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("*** Deck of %d cards ***\n", d.GetLength()))
	out.WriteString(fmt.Sprintf("Deck count: %d\n", d.DeckCount))
	out.WriteString(fmt.Sprintf("Shuffle count: %d\n", d.ShuffleCount))
	out.WriteString("{\n")
	for _, card := range d.Cards {
		out.WriteString(fmt.Sprintf("\t%s\n", card.Debug()))
	}
	out.WriteString("}\n")

	return out.String()
}

// TODO: alter blackjackdeck so that it always maintains slice of original
// another slice for the deck that is going to shrink after cards of taken
// and then every card 'delt' is also a slice, but they point to the second slice
// then when a reshuffle happens cause cards are depleted
