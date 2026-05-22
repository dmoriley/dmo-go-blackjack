// Package suit provides types and helpers for card suits
package suit

import "fmt"

const (
	Diamonds = "Diamonds"
	Hearts   = "Hearts"
	Spades   = "Spades"
	Clubs    = "Clubs"
)

var Suits = map[string]struct{}{
	Diamonds: {},
	Hearts:   {},
	Spades:   {},
	Clubs:    {},
}

func NewSuit(suit string) (string, error) {
	if _, ok := Suits[suit]; !ok {
		return "", fmt.Errorf("%q is not a valid card suit", suit)
	}

	return suit, nil
}
