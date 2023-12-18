package deck

import (
	"bytes"
	"fmt"
	"blackjack/card"
	"blackjack/card/rank"
	"blackjack/card/suit"
	"math/rand"
)

func new(ranks map[string]int) *Deck {
	deck := &Deck{}
	cards := []*card.Card{}

	for cardSuit := range suit.Suits {
		for cardRank, cardValue := range ranks {
			card, error := card.NewCard(cardSuit, cardRank, cardValue)
			if error != nil {
				panic(error)
			}

			cards = append(cards, card)
		}
	}
	deck.cards = cards
	return deck
}

func NewDefaultDeck() *Deck {
	return new(rank.Ranks)
}

func NewCustomDeck(ranks map[string]int) *Deck {
	return new(ranks)
}

// ************* Deck type ***************

type Deck struct {
	cards        []*card.Card
	shuffleCount int
}

func (d *Deck) getLength() int {
	return len(d.cards)
}

func (d *Deck) Inspect() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("*** Deck of %d cards ***\n", d.getLength()))
	out.WriteString(fmt.Sprintf("Shuffle count: %d\n", d.shuffleCount))
	out.WriteString("{\n")
	for _, card := range d.cards {
		out.WriteString(fmt.Sprintf("\t%s\n", card.Inspect()))
	}
	out.WriteString("}\n")

	return out.String()
}

func (d *Deck) Shuffle(count int) {
	currentIndex := d.getLength()
	var randomIndex int

	for currentIndex > 0 {
		randomIndex = rand.Intn(d.getLength())
		currentIndex--

		// swap  current index with random index
		d.cards[currentIndex], d.cards[randomIndex] = d.cards[randomIndex], d.cards[currentIndex]
	}

	d.shuffleCount++

	// recursively shuffle
	if count > 1 {
		d.Shuffle(count - 1)
	}
}
