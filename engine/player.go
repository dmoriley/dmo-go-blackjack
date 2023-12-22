package engine

import (
	"blackjack/card"
)

type Player struct {
	Name  string
	Cash  int
	Cards []*card.Card
}
