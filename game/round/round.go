package round

import "blackjack/card"

type Dealer struct {
	cards *[]card.Card
}

type Player struct {
	Cards *[]card.Card
	Bet   int
}

type Round struct {
	Player Player
	Dealer Dealer
}
