package card

import (
	"bytes"
	"fmt"
	"blackjack/card/rank"
	"blackjack/card/suit"
)

func NewCard(suitInput string, rankNameInput string, cardValue int) (*Card, error) {
	cardRank, error := rank.NewRank(rankNameInput, cardValue)

	if error != nil {
		return nil, error
	}

	cardSuit, error := suit.NewSuit(suitInput)

	if error != nil {
		return nil, error
	}

	card := &Card{
		Suit:     cardSuit,
		Rank:     cardRank,
		IsFaceUp: false,
	}

	return card, nil
}

type Card struct {
	Suit     string
	Rank     *rank.Rank
	IsFaceUp bool
}

func (c *Card) Inspect() string {
	var out bytes.Buffer

	out.WriteString("{")
	out.WriteString(fmt.Sprintf("%s of %s", c.Rank.Name, c.Suit))
	out.WriteString(fmt.Sprintf(", value: %d", c.Rank.Value))
	out.WriteString("}")

	return out.String()
}
