package deck

import (
	"blackjack/card"
	"blackjack/card/rank"
	"blackjack/card/suit"
	"bytes"
	"fmt"
	"math/rand"
)

func newDeck(ranks map[string]int) *Deck {
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
	deck.Cards = cards
	return deck
}

func NewDefaultDeck() *Deck {
	return newDeck(rank.Ranks)
}

func NewCustomDeck(ranks map[string]int) *Deck {
	return newDeck(ranks)
}

// ************* Deck type ***************

type Deck struct {
	Cards        []*card.Card
	ShuffleCount int
}

func (d *Deck) GetLength() int {
	return len(d.Cards)
}

func (d *Deck) Inspect() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("*** Deck of %d cards ***\n", d.GetLength()))
	out.WriteString(fmt.Sprintf("Shuffle count: %d\n", d.ShuffleCount))
	out.WriteString(PrintCards(d.Cards))

	return out.String()
}

func (d *Deck) Shuffle(count int) {
	currentIndex := d.GetLength()
	var randomIndex int

	for currentIndex > 0 {
		randomIndex = rand.Intn(d.GetLength())
		currentIndex--

		// swap  current index with random index
		d.Cards[currentIndex], d.Cards[randomIndex] = d.Cards[randomIndex], d.Cards[currentIndex]
	}

	d.ShuffleCount++

	// recursively shuffle
	if count > 1 {
		d.Shuffle(count - 1)
	}
}

func (d *Deck) AddDeck(deckToAdd *Deck) {
	d.Cards = append(d.Cards, deckToAdd.Cards...)
}

func (d *Deck) AddCard(cardToAdd *card.Card) {
	d.Cards = append(d.Cards, cardToAdd)
}

func PrintCards(cards []*card.Card) string {
	var out bytes.Buffer

	out.WriteString("{\n")
	for _, card := range cards {
		out.WriteString(fmt.Sprintf("\t%s\n", card.Inspect()))
	}
	out.WriteString("}\n")

	return out.String()
}
