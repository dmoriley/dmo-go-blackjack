package rank

import (
	"fmt"
)

const (
	Ace   = "Ace"
	Two   = "Two"
	Three = "Three"
	Four  = "Four"
	Five  = "Five"
	Six   = "Six"
	Seven = "Seven"
	Eight = "Eight"
	Nine  = "Nine"
	Ten   = "Ten"
	Jack  = "Jack"
	Queen = "Queen"
	King  = "King"
)

var Ranks = map[string]int{
	Ace:   1,
	Two:   2,
	Three: 3,
	Four:  4,
	Five:  5,
	Six:   6,
	Seven: 7,
	Eight: 8,
	Nine:  9,
	Ten:   10,
	Jack:  11,
	Queen: 12,
	King:  13,
}

// ************** Rank ***************

type Rank struct {
	Name  string
	Value int
}

func NewRank(rankName string, value int) (*Rank, error) {
	if _, ok := Ranks[rankName]; !ok {
		return nil, fmt.Errorf("%q is not a valid card rank", rankName)
	}

	return &Rank{
		Name:  rankName,
		Value: value,
	}, nil
}

func (r *Rank) Inspect() string {
	return fmt.Sprintf("{ rank: %s, value: %d }", r.Name, r.Value)
}
