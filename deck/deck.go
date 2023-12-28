package deck

import (
	"blackjack/card"
	"blackjack/card/rank"
	"blackjack/card/suit"
	"bytes"
	"fmt"
	"math/rand"
)

func newDeck(suits []string, ranks map[string]int) *Deck {
	deck := &Deck{}
	cards := []*card.Card{}

	for _, cardSuit := range suits {
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
	return newDeck([]string{suit.Hearts, suit.Clubs, suit.Diamonds, suit.Spades}, rank.Ranks)
}

func NewCustomDeck(suits []string, ranks map[string]int) *Deck {
	return newDeck(suits, ranks)
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

// Add a list of cards to this deck's cards
func (d *Deck) AddCards(cards []*card.Card) {
	d.Cards = append(d.Cards, cards...)
}

// Print and format a list of cards
func PrintCards(cards []*card.Card) string {
	var out bytes.Buffer

	out.WriteString("{\n")
	for _, card := range cards {
		out.WriteString(fmt.Sprintf("\t%s\n", card.Inspect()))
	}
	out.WriteString("}\n")

	return out.String()
}

func (d *Deck) Pop(count int) (popped []*card.Card) {
	if count == 0 {
		count = 1
	}

	if len(d.Cards) == 0 || count-1 > len(d.Cards) {
		return nil
	}

	// get slice of cards from beginning to count inclusive
	// and copy the contents to a new slice so as not to hold
	// onto the memory of the original for potential memory leak
	copy(popped, d.Cards[:count-1])
	// re-slice the original array to start after the popped
	d.Cards = d.Cards[count-1:]

	return
}
