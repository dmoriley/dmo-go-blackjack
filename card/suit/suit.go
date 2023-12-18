package suit

import "fmt"

const (
	Diamonds = "Diamonds"
	Hearts   = "Hearts"
	Spades   = "Spades"
	Clubs    = "Clubs"
)

var Suits = map[string]bool{
	Diamonds: true,
	Hearts:   true,
	Spades:   true,
	Clubs:    true,
}

func NewSuit(suit string) (string, error) {
	if _, ok := Suits[suit]; !ok {
		return "", fmt.Errorf("%q is not a valid card suit", suit)
	}

	return suit, nil
}
