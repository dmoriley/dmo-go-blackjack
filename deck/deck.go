package deck

import (
	"blackjack/card"
	"blackjack/card/rank"
	"blackjack/card/suit"
	"bytes"
	"fmt"
	"math/rand"
)

type Deck struct {
	Cards        []*card.Card
	ShuffleCount int
}

func newDeck(suits []string, ranks map[string]int) *Deck {
	Deck := &Deck{}
	cards := []*card.Card{}

	for _, cardSuit := range suits {
		for cardRank, cardValue := range ranks {
			card, error := card.NewCard(cardSuit, cardRank, cardValue, false)
			if error != nil {
				panic(error)
			}

			cards = append(cards, card)
		}
	}
	Deck.Cards = cards
	return Deck
}

type Opts struct {
	suits []string
	ranks map[string]int
}

type OptsFunc func(*Opts)

func defaultOps() Opts {
	return Opts{
		suits: []string{suit.Hearts, suit.Clubs, suit.Diamonds, suit.Spades},
		ranks: rank.Ranks,
	}
}

func WithSuits(suits []string) OptsFunc {
	return func(opts *Opts) {
		opts.suits = suits
	}
}

func WithRanks(ranks map[string]int) OptsFunc {
	return func(opts *Opts) {
		opts.ranks = ranks
	}
}

// Create a new Deck. Defaults to a regular 52 card deck. Optionally
// set custom suits or ranks using the WithSuits or WithRanks optional function
// parameters
func NewDeck(opts ...OptsFunc) *Deck {
	// optional method parameters using the "functional options" pattern
	o := defaultOps()

	for _, fn := range opts {
		fn(&o)
	}
	return newDeck(o.suits, o.ranks)
}

func (d *Deck) GetLength() int {
	return len(d.Cards)
}

func (d *Deck) Inspect() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("*** Deck of %d cards ***\n", d.GetLength()))
	out.WriteString(fmt.Sprintf("Shuffle count: %d\n", d.ShuffleCount))

	return out.String()
}

func (d *Deck) Debug() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("*** Deck of %d cards ***\n", d.GetLength()))
	out.WriteString(fmt.Sprintf("Shuffle count: %d\n", d.ShuffleCount))
	out.WriteString(PrintCards(d.Cards, true))

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

// Get a card(s) from the deck
func (d *Deck) Pop(count int) []*card.Card {
	if count == 0 {
		count = 1
	}

	if len(d.Cards) == 0 || count-1 > len(d.Cards) {
		return nil
	}

	var popped = make([]*card.Card, count)
	// get slice of cards from beginning to count inclusive
	// and copy the contents to a new slice so as not to hold
	// onto the memory of the original for potential memory leak
	copy(popped, d.Cards[:count])
	// re-slice the original array to start after the popped
	d.Cards = d.Cards[count:]

	return popped
}

// Print and format a list of cards
func PrintCards(cards []*card.Card, printDebug bool) string {
	var out bytes.Buffer

	if len(cards) != 0 {
		out.WriteString("{\n")
		if printDebug == false {
			for _, card := range cards {
				out.WriteString(fmt.Sprintf("\t%s\n", card.Inspect()))
			}
		} else {
			for _, card := range cards {
				out.WriteString(fmt.Sprintf("\t%s\n", card.Debug()))
			}
		}
		out.WriteString("}\n")
	} else {
		out.WriteString("{ No Cards }")
	}

	return out.String()
}

func PrettyPrintCards(cards []*card.Card) string {
	// TODO: pretty print ASCII looking cards
	// .------..------..------..------.     .------..------..------..------..------..------..------.
	// |M.--. ||O.--. ||R.--. ||E.--. |.-.  |O.--. ||P.--. ||T.--. ||I.--. ||O.--. ||N.--. ||S.--. |
	// | (\/) || :/\: || :(): || (\/) ((5)) | :/\: || :/\: || :/\: || (\/) || :/\: || :(): || :/\: |
	// | :\/: || :\/: || ()() || :\/: |'-.-.| :\/: || (__) || (__) || :\/: || :\/: || ()() || :\/: |
	// | '--'M|| '--'O|| '--'R|| '--'E| ((1)) '--'O|| '--'P|| '--'T|| '--'I|| '--'O|| '--'N|| '--'S|
	// `------'`------'`------'`------'  '-'`------'`------'`------'`------'`------'`------'`------'
	return PrintCards(cards, false)
}
