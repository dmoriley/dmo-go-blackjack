package card

import (
	"blackjack/card/rank"
	"blackjack/card/suit"
	"bytes"
	"fmt"
)

func NewCard(suitInput string, rankNameInput string, cardValue int, isFaceUp bool) (*Card, error) {
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
		IsFaceUp: isFaceUp,
	}

	return card, nil
}

type Card struct {
	Suit     string
	Rank     *rank.Rank
	IsFaceUp bool
}

func (c *Card) Debug() string {
	var out bytes.Buffer

	out.WriteString("{")
	out.WriteString(fmt.Sprintf("%s of %s", c.Rank.Name, c.Suit))
	out.WriteString(fmt.Sprintf(", value: %d", c.Rank.Value))
	out.WriteString(fmt.Sprintf(", IsFaceUp: %t", c.IsFaceUp))
	out.WriteString("}")

	return out.String()
}

func (c *Card) Inspect() string {
	var out bytes.Buffer

	if c.IsFaceUp {
		out.WriteString("{")
		out.WriteString(fmt.Sprintf("%s of %s", c.Rank.Name, c.Suit))
		out.WriteString(fmt.Sprintf(", value: %d", c.Rank.Value))
		out.WriteString("}")
	} else {
		out.WriteString("{")
		out.WriteString("Face down")
		out.WriteString("}")
	}

	return out.String()
}
