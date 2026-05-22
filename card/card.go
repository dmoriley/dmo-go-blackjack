// Package card provides types and helpers for playing cards
package card

import (
	"bytes"
	"fmt"

	"blackjack/card/rank"
	"blackjack/card/suit"
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

func (c Card) Debug() string {
	var out bytes.Buffer

	out.WriteString("{")
	fmt.Fprintf(&out, "%s of %s", c.Rank.Name, c.Suit)
	fmt.Fprintf(&out, ", value: %d", c.Rank.Value)
	fmt.Fprintf(&out, ", IsFaceUp: %t", c.IsFaceUp)
	out.WriteString("}")

	return out.String()
}

func (c Card) Inspect() string {
	var out bytes.Buffer

	if c.IsFaceUp {
		out.WriteString("{")
		fmt.Fprintf(&out, "%s of %s", c.Rank.Name, c.Suit)
		fmt.Fprintf(&out, ", value: %d", c.Rank.Value)
		out.WriteString("}")
	} else {
		out.WriteString("{")
		out.WriteString("Face down")
		out.WriteString("}")
	}

	return out.String()
}
