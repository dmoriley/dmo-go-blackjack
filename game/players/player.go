package players

import "blackjack/card"

type Player struct {
	Cards       []*card.Card
	Name        string
	Cash        int
	Bet         int
	PreviousBet int
	SplitCards  []*card.Card
}

// Move all cards to split cards
func (p *Player) MoveCardsToSplit() {
	// make split cards new slice of cards
	p.SplitCards = p.Cards[:]
	// reslice cards just past length so its effectively empty
	p.Cards = p.Cards[len(p.Cards):]
}

func (p *Player) NextSplitCard() {
	if len(p.SplitCards) == 0 {
		return
	}

	p.Cards = append(p.Cards, p.SplitCards[0])
	// reslice splitCards just past card that was appended
	p.SplitCards = p.SplitCards[1:]
}

func (p *Player) HasSplitCards() bool {
	return len(p.SplitCards) > 0
}
