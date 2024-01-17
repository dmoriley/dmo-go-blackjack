package decks

import (
	"blackjack/card"
	"blackjack/card/rank"
	"bytes"
	"fmt"
)

type BlackjackDeck struct {
	// Number of decks in the blackjack deck
	DeckCount int
	// When the deck reaches a minimum number of cards reshuffle all cards back to deck
	minCardCount int
	// Slice of all the cards before any are removed
	// TODO: possibly remove this property
	allCards []*card.Card
	// Cards that have been used and are no longer in play
	discardedCards []*card.Card
	Deck
}

type BlackjackDeckConfig struct {
	numberOfDecks int
	minCardCount  int
}

func NewBlackjackDeckConfig() *BlackjackDeckConfig {
	return &BlackjackDeckConfig{
		// default values
		numberOfDecks: 6,
		// minCardCount:  60,
		minCardCount: 6*52 - 10,
	}
}

// Number of decks the blackjack deck should contain
func (c *BlackjackDeckConfig) WithNumberOfDecks(count int) *BlackjackDeckConfig {
	c.numberOfDecks = count
	return c
}

// Minimum number of cards before the deck is reshuffled
func (c *BlackjackDeckConfig) WithMinCardCount(count int) *BlackjackDeckConfig {
	c.minCardCount = count
	return c
}

func NewBlackjackDeck(config *BlackjackDeckConfig) *BlackjackDeck {
	bjDeck := &BlackjackDeck{
		DeckCount:    0,
		minCardCount: config.minCardCount,
	}

	deckToAdd := NewDeck()

	// change all the kings, queens and jacks to value of 10
	for _, card := range deckToAdd.Cards {
		if card.Rank.Name == rank.King || card.Rank.Name == rank.Queen ||
			card.Rank.Name == rank.Jack {
			card.Rank.Value = 10
		}
	}

	for i := 0; i < config.numberOfDecks; i++ {
		bjDeck.AddCards(deckToAdd.Cards)
		bjDeck.DeckCount++
	}

	bjDeck.allCards = bjDeck.Cards[:] // get a slice of all the cards

	fmt.Printf("\nCards len: %d, Cards cap: %d\n", bjDeck.GetLength(), cap(bjDeck.Cards))
	return bjDeck
}

func (d *BlackjackDeck) Inspect() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("*** Deck of %d cards ***\n", d.GetLength()))
	out.WriteString(fmt.Sprintf("Deck count: %d\n", d.DeckCount))
	out.WriteString(fmt.Sprintf("Shuffle count: %d\n", d.ShuffleCount))
	out.WriteString(fmt.Sprintf("Minimum card count: %d\n", d.minCardCount))

	return out.String()
}

func (d *BlackjackDeck) Debug() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("*** Deck of %d cards ***\n", d.GetLength()))
	out.WriteString(fmt.Sprintf("Deck count: %d\n", d.DeckCount))
	out.WriteString(fmt.Sprintf("Shuffle count: %d\n", d.ShuffleCount))
	out.WriteString(fmt.Sprintf("Minimum card count: %d\n", d.minCardCount))
	out.WriteString("{\n")
	for _, card := range d.Cards {
		out.WriteString(fmt.Sprintf("\t%s\n", card.Debug()))
	}
	out.WriteString("}\n")

	return out.String()
}

func (d *BlackjackDeck) Pop(count int) []*card.Card {
	if count == 0 {
		count = 1
	}

	if d.GetLength() == 0 || count-1 > len(d.Cards) {
		return nil
	}

	var popped = make([]*card.Card, count)
	// get slice of cards from beginning to count inclusive
	// and copy the contents to a new slice so as not to hold
	// onto the memory of the original for potential memory leak
	copy(popped, d.Cards[:count])
	// re-slice the original array to start after the popped
	d.Cards = d.Cards[count:]

	// TODO: need to rethink this logic cause it duplicates the cards that are still in play on the table
	if d.GetLength() == d.minCardCount {
		d.Reshuffle(5)
	}

	return popped

}

// Added cards to the discarded pile
func (d *BlackjackDeck) AddDiscardedCards(cards []*card.Card) {
	d.discardedCards = append(d.discardedCards, cards...)
}

// Added the discarded cards back to the deck and shuffle
func (d *BlackjackDeck) Reshuffle(shuffleCount int) {
	fmt.Println("\n***********************************")
	fmt.Println("***** Reshuffling the deck... *****")
	fmt.Println("***********************************")
	d.Cards = append(d.Cards, d.discardedCards...)
	d.Shuffle(shuffleCount)
}
