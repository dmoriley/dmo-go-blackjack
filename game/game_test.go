package game

import (
	"blackjack/card"
	"blackjack/card/rank"
	"blackjack/card/suit"
	"blackjack/decks"
	"testing"
)

func TestTenCards(t *testing.T) {
	jack, _ := card.NewCard(suit.Hearts, rank.Jack, 10, true)
	queen, _ := card.NewCard(suit.Spades, rank.Queen, 10, true)
	king, _ := card.NewCard(suit.Diamonds, rank.King, 10, true)

	cards := []*card.Card{
		jack,
	}
	want := 10
	got := GetCardsTotal(cards)

	if want != got {
		t.Errorf("Card total wrong. Want = %d, got = %d", want, got)
	}

	cards = append(cards, queen)
	want = 20
	got = GetCardsTotal(cards)

	if want != got {
		t.Errorf("Card total wrong. Want = %d, got = %d", want, got)
	}

	cards = append(cards, king)
	want = 30
	got = GetCardsTotal(cards)

	if want != got {
		t.Errorf("Card total wrong. Want = %d, got = %d", want, got)
	}

}

func TestNumberCards(t *testing.T) {
	cards := []*card.Card{}
	two, _ := card.NewCard(suit.Hearts, rank.Two, 2, true)
	three, _ := card.NewCard(suit.Hearts, rank.Three, 3, true)
	four, _ := card.NewCard(suit.Hearts, rank.Four, 4, true)
	five, _ := card.NewCard(suit.Hearts, rank.Five, 5, true)
	six, _ := card.NewCard(suit.Hearts, rank.Six, 6, true)
	seven, _ := card.NewCard(suit.Hearts, rank.Seven, 7, true)
	eight, _ := card.NewCard(suit.Hearts, rank.Eight, 8, true)
	nine, _ := card.NewCard(suit.Hearts, rank.Nine, 9, true)

	cards = append(cards, two)
	want := 2
	got := GetCardsTotal(cards)
	if want != got {
		t.Errorf("Card total wrong. Want = %d, got = %d", want, got)
	}

	cards = append(cards, three)
	want += 3
	got = GetCardsTotal(cards)
	if want != got {
		t.Errorf("Card total wrong. Want = %d, got = %d", want, got)
	}

	cards = append(cards, four)
	want += 4
	got = GetCardsTotal(cards)
	if want != got {
		t.Errorf("Card total wrong. Want = %d, got = %d", want, got)
	}

	cards = append(cards, five)
	want += 5
	got = GetCardsTotal(cards)
	if want != got {
		t.Errorf("Card total wrong. Want = %d, got = %d", want, got)
	}

	cards = append(cards, six)
	want += 6
	got = GetCardsTotal(cards)
	if want != got {
		t.Errorf("Card total wrong. Want = %d, got = %d", want, got)
	}

	cards = append(cards, seven)
	want += 7
	got = GetCardsTotal(cards)
	if want != got {
		t.Errorf("Card total wrong. Want = %d, got = %d", want, got)
	}

	cards = append(cards, eight)
	want += 8
	got = GetCardsTotal(cards)
	if want != got {
		t.Errorf("Card total wrong. Want = %d, got = %d", want, got)
	}

	cards = append(cards, nine)
	want += 9
	got = GetCardsTotal(cards)
	if want != got {
		t.Errorf("Card total wrong. Want = %d, got = %d", want, got)
	}
}

func TestAceCard(t *testing.T) {
	ace, _ := card.NewCard(suit.Hearts, rank.Ace, 1, true)
	cards := []*card.Card{
		ace,
	}

	want := 11
	got := GetCardsTotal(cards)
	if want != got {
		t.Errorf("Card total wrong. Want = %d, got = %d", want, got)
	}

	cards = append(cards, ace)
	want = 12
	got = GetCardsTotal(cards)
	if want != got {
		t.Errorf("Card total wrong. Want = %d, got = %d", want, got)
	}

	// increment it past 21
	for i := 0; i < 9; i++ {
		cards = append(cards, ace)
		want++
	}

	got = GetCardsTotal(cards)
	// want is 21
	if want != got {
		t.Errorf("Card total wrong. Want = %d, got = %d", want, got)
	}

	// adding one more ace should make that 12 aces, so all aces at that point should be value of one
	// cause one of them being 11 would push it over bust limit
	cards = append(cards, ace)
	want = 12
	got = GetCardsTotal(cards)
	if want != got {
		t.Errorf("Card total wrong. Want = %d, got = %d", want, got)
	}
}

func TestCombiningAllCards(t *testing.T) {
	cards := []*card.Card{}
	ace, _ := card.NewCard(suit.Hearts, rank.Ace, 1, true)
	two, _ := card.NewCard(suit.Hearts, rank.Two, 2, true)
	three, _ := card.NewCard(suit.Hearts, rank.Three, 3, true)
	// four, _ := card.NewCard(suit.Hearts, rank.Four, 4, true)
	five, _ := card.NewCard(suit.Hearts, rank.Five, 5, true)
	six, _ := card.NewCard(suit.Hearts, rank.Six, 6, true)
	// seven, _ := card.NewCard(suit.Hearts, rank.Seven, 7, true)
	// eight, _ := card.NewCard(suit.Hearts, rank.Eight, 8, true)
	// nine, _ := card.NewCard(suit.Hearts, rank.Nine, 9, true)
	jack, _ := card.NewCard(suit.Hearts, rank.Jack, 10, true)
	queen, _ := card.NewCard(suit.Spades, rank.Queen, 10, true)
	king, _ := card.NewCard(suit.Diamonds, rank.King, 10, true)

	// ace with ten card should be 21
	cards = append(cards, ace)
	cards = append(cards, king)

	want := 21
	got := GetCardsTotal(cards)
	if want != got {
		t.Errorf("Card total wrong. Want = %d, got = %d", want, got)
	}

	// adding another ace should change the value of both ace's to 1's
	cards = append(cards, ace)
	want = 10 + 1 + 1 // 12
	got = GetCardsTotal(cards)
	if want != got {
		t.Errorf("Card total wrong. Want = %d, got = %d", want, got)
	}

	cards = []*card.Card{
		queen, six, ace,
	}

	want = 10 + 6 + 1 // 17
	got = GetCardsTotal(cards)
	if want != got {
		t.Errorf(
			"Card total wrong. Want = %d, got = %d\n%s",
			want,
			got,
			decks.PrintCards(cards, true),
		)

	}

	cards = []*card.Card{
		two, three, five, jack, ace,
	}

	want = 2 + 3 + 5 + 10 + 1 // 21
	got = GetCardsTotal(cards)
	if want != got {
		t.Errorf(
			"Card total wrong. Want = %d, got = %d\n%s",
			want,
			got,
			decks.PrintCards(cards, true),
		)

	}

	cards = []*card.Card{
		king, queen, ace,
	}

	want = 10 + 10 + 1
	got = GetCardsTotal(cards)
	if want != got {
		t.Errorf(
			"Card total wrong. Want = %d, got = %d\n%s",
			want,
			got,
			decks.PrintCards(cards, true),
		)

	}
}
