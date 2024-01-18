package decks

import (
	"blackjack/card"
	"testing"
)

// Test that every card in the blackjack deck is unique
// and that there arent multiple pointers to the same card
func TestBlackjackDeckOnlyUniqueCards(t *testing.T) {
	config := NewBlackjackDeckConfig()
	deck := NewBlackjackDeck(config)
	originalLength := deck.GetLength()

	allFaceDown := true
	for _, card := range deck.Cards {
		if card.IsFaceUp {
			// found card not face down, thats not right
			allFaceDown = false
			break
		}
	}

	if !allFaceDown {
		t.Fatalf("Not all cards are face down after deck initialization")
	}

	// pull out one poppedCard and change it to face up
	poppedCard := deck.Pop(1)[0]
	poppedCard.IsFaceUp = true
	newLength := deck.GetLength()

	if originalLength-1 != newLength {
		t.Fatalf(
			"New length doesn't make any sense. Should be original length minus one but got %d",
			newLength,
		)
	}

	// check if all cards left in the deck are face down
	// if one is face up, means it was as pointer to the same card
	faceUpCards := []*card.Card{}
	for _, card := range deck.Cards {
		if card.IsFaceUp {
			// found card not face down, thats not right
			faceUpCards = append(faceUpCards, card)
		}
	}

	if len(faceUpCards) != 0 {
		t.Fatalf(
			"Founds card(s) that are face up when they shouldn't be. Count: %d\n%s",
			len(faceUpCards),
			PrintCards(faceUpCards, true),
		)
	}

}

// Test that all cards are returned from the discard pile and that
// all cards are flipped back to face down on a reshuffle
func TestBlackjackReshuffle(t *testing.T) {
	config := NewBlackjackDeckConfig()
	deck := NewBlackjackDeck(config)
	originalLength := deck.GetLength()

	allFaceDown := true
	for _, card := range deck.Cards {
		if card.IsFaceUp {
			// found card not face down, thats not right
			allFaceDown = false
			break
		}
	}

	if !allFaceDown {
		t.Fatalf("Not all cards are face down after deck initialization")
	}

	// pull out 5 cards and put in the discard pile
	for i := 0; i < 5; i++ {
		c := deck.Pop(1)
		c[0].IsFaceUp = true
		deck.AddDiscardedCards(c)
	}

	if len(deck.discardedCards) != 5 {
		t.Fatalf(
			"Discarded cards not the right length. Expected %d but got %d",
			5,
			len(deck.discardedCards),
		)
	}

	if len(deck.Cards) != originalLength-5 {
		t.Fatalf(
			"deck cards length is wrong. Expected %d and got %d",
			originalLength-5,
			len(deck.Cards),
		)
	}

	deck.Reshuffle(5)

	if len(deck.discardedCards) != 0 {
		t.Fatalf(
			"Discarded cards should be empty after a reshuffle. Instead got %d",
			len(deck.discardedCards),
		)
	}
	// check if all cards left in the deck are face down
	// if one is face up, means it was as pointer to the same card
	faceUpCards := []*card.Card{}
	for _, card := range deck.Cards {
		if card.IsFaceUp {
			// found card not face down, thats not right
			faceUpCards = append(faceUpCards, card)
		}
	}

	if len(faceUpCards) != 0 {
		t.Fatalf(
			"Founds card(s) that are face up when there should be zero. Count: %d\n%s",
			len(faceUpCards),
			PrintCards(faceUpCards, true),
		)
	}

}
