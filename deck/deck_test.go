package deck

import (
	"blackjack/card/rank"
	"blackjack/card/suit"
	"testing"
)

func TestNewDeckWithoutOptions(t *testing.T) {
	// should create a default deck of 52 cards
	deck := NewDeck()
	dlen := len(deck.Cards)
	if dlen != 52 {
		t.Fatalf("Expected %d cards in deck, got %d", 52, dlen)
	}

}

func TestCustomSuitsNewDeck(t *testing.T) {
	suits := []string{suit.Diamonds, suit.Clubs}
	expectedDeckLength := len(suits) * len(rank.Ranks)
	deck := NewDeck(WithSuits(suits))
	dlen := len(deck.Cards)

	if dlen != expectedDeckLength {
		t.Fatalf("Expected %d cards in deck, got %d", expectedDeckLength, dlen)
	}
}

func TestCustomRanksNewDeck(t *testing.T) {
	ranks := map[string]int{
		"Ace":  100,
		"King": 50,
	}
	expectedDeckLength := len(suit.Suits) * len(ranks)
	deck := NewDeck(WithRanks(ranks))
	dlen := len(deck.Cards)

	if dlen != expectedDeckLength {
		t.Fatalf("Expected %d cards in deck, got %d", expectedDeckLength, dlen)
	}
}

// TODO: finish testing these later
/* func TestAddCards(t *testing.T) {

}

func TestPop(t *testing.T) {

} */
