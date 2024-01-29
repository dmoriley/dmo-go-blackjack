package players

import (
	"blackjack/card"
	"blackjack/card/rank"
	"blackjack/card/suit"
	"testing"
)

func TestNextSplitCardMovesCardCorrectly(t *testing.T) {
	t.Skip()
	ace, _ := card.NewCard(suit.Hearts, rank.Ace, 1, true)
	two, _ := card.NewCard(suit.Hearts, rank.Two, 2, true)
	three, _ := card.NewCard(suit.Hearts, rank.Three, 3, true)
	four, _ := card.NewCard(suit.Hearts, rank.Four, 4, true)

	player := &Player{
		Cards: []*card.Card{
			ace, two,
		},
		Name: "Tester",
		Cash: 0,
		Bet:  0,
		SplitCards: []*card.Card{
			three, four,
		},
	}

	player.NextSplitCard()

	if len(player.SplitCards) != 1 && player.SplitCards[0] != four {
		t.Fatalf(
			"Split cards are wrong. Expected length of %d but got %d. Card is %s",
			1,
			len(player.SplitCards),
			player.SplitCards[0].Inspect(),
		)
	}

	if len(player.Cards) != 3 {
		t.Fatalf("Player cards length is wrong. Expected %d but got %d", 3, len(player.Cards))
	}

	if player.Cards[2] != three {
		t.Fatalf(
			"Wrong card was added to player cards. Wanted %q but got %q",
			three.Inspect(),
			player.Cards[2].Inspect(),
		)
	}
}

func TestNextSplitCardWithOnlyOneCardLeft(t *testing.T) {
	ace, _ := card.NewCard(suit.Hearts, rank.Ace, 1, true)
	two, _ := card.NewCard(suit.Hearts, rank.Two, 2, true)
	three, _ := card.NewCard(suit.Hearts, rank.Three, 3, true)
	four, _ := card.NewCard(suit.Hearts, rank.Four, 4, true)

	player := &Player{
		Cards: []*card.Card{
			ace, two, three,
		},
		Name: "Tester",
		Cash: 0,
		Bet:  0,
		SplitCards: []*card.Card{
			four,
		},
	}

	player.NextSplitCard()

	if len(player.SplitCards) != 0 {
		t.Fatalf(
			"Player split cards wrong. Should be empty but got length of %d",
			len(player.SplitCards),
		)
	}
}
