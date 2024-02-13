package utils

import (
	"blackjack/card"
	"blackjack/card/rank"
)

const BLACKJACK = 21

// Determine the card value total of the hand supplied
// Does not count value of card that is not face up
func CalcCardsTotal(cards []*card.Card) int {
	total := 0
	aceCount := 0

	for _, card := range cards {
		if !card.IsFaceUp {
			// dont total cards that arent being shown
			continue
		}
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
