package engine

import (
	"blackjack/card"
	"blackjack/card/rank"
	"blackjack/deck"
	"blackjack/engine"
)

const (
	// Max value to win or bust
	BLACKJACK = 21
)

type Blackjack struct {
	Player *engine.Player
	Dealer *engine.Dealer
	Deck   *deck.Deck
}

// Determine the card value total of the hand supplied
func GetCardsTotal(cards []*card.Card) int {
	total := 0
	aceCount := 0

	for _, card := range cards {
		switch card.Rank.Name {
		case rank.Ace:
			aceCount++
		default:
			total += card.Rank.Value
		}
	}

	if aceCount == 1 {
		if total+11 <= BLACKJACK {
			total += 11
		} else {
			total += 1
		}
	} else if aceCount > 1 {
		// check if one ace can have the value of 11
		if total+11+(aceCount-1) <= BLACKJACK {
			total += (11 + (aceCount - 1))
		} else {
			// all aces must have a value of one
			total += aceCount
		}
	}

	return total
}
