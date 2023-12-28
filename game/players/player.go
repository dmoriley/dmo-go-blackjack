package players

import "blackjack/card"

type Player struct {
	Cards []*card.Card
	Name  string
	Cash  int
	Bet   int
}
