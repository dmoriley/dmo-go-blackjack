package game

import (
	"blackjack/card"
	"blackjack/card/rank"
	"blackjack/card/suit"
	"blackjack/decks"
	"blackjack/game/players"
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

func TestDealerSoft17(t *testing.T) {

	jack, _ := card.NewCard(suit.Hearts, rank.Jack, 2, true)
	king, _ := card.NewCard(suit.Hearts, rank.King, 10, true)

	// total 20
	playerCards := []*card.Card{
		jack,
		king,
	}

	ace, _ := card.NewCard(suit.Hearts, rank.Ace, 1, true)
	six, _ := card.NewCard(suit.Hearts, rank.Six, 6, true)

	// should be a soft 17
	dealerCards := []*card.Card{
		ace,
		six,
	}

	if dt := GetCardsTotal(dealerCards); dt != 17 {
		t.Fatalf("Card total not 17, got %d", dt)
	}

	five, _ := card.NewCard(suit.Hearts, rank.Five, 5, true)
	four, _ := card.NewCard(suit.Hearts, rank.Four, 4, true)

	bj := &Blackjack{
		Player: &players.Player{
			Cards: playerCards,
			Name:  "player1",
			Cash:  500,
			Bet:   5,
		},
		Dealer: &players.Dealer{
			Cards: dealerCards,
		},
		// should have minCardCount of 0 as 'zero' value for being unset
		Deck: &decks.BlackjackDeck{
			DeckCount: 1,
			Deck: decks.Deck{
				Cards: []*card.Card{
					// first card out of the deck should be six
					six,
					five,
					four,
				},
			},
		},
	}

	outcome := bj.PlayerStand()

	if len(bj.Dealer.Cards) != 4 {
		t.Fatalf("Dealer doesn't have 4 cards, got %d", len(bj.Dealer.Cards))
	}

	expectedTotal := 1 + six.Rank.Value + six.Rank.Value + five.Rank.Value
	if dt := GetCardsTotal(bj.Dealer.Cards); dt != expectedTotal {
		t.Fatalf("Card total not %d, got %d", expectedTotal, dt)
	}

	if outcome != PlayerLost {
		t.Fatalf("Round outcome is wrong. Should have %d, but got %d", PlayerLost, outcome)
	}
}

func TestDoubleMoveAvailable(t *testing.T) {
	// double is allowed when a intial two cards has a total of 9-11

	four, _ := card.NewCard(suit.Hearts, rank.Four, 4, true)
	six, _ := card.NewCard(suit.Hearts, rank.Six, 6, true)

	// total of 10
	playerCards := []*card.Card{
		four, six,
	}

	five, _ := card.NewCard(suit.Hearts, rank.Five, 5, true)
	seven, _ := card.NewCard(suit.Hearts, rank.Seven, 7, true)

	// total of 12
	dealerCards := []*card.Card{
		five,
		seven,
	}

	bj := &Blackjack{
		Player: &players.Player{
			Cards: playerCards,
			Name:  "player1",
			Cash:  500,
			Bet:   5,
		},
		Dealer: &players.Dealer{
			Cards: dealerCards,
		},
		// should have minCardCount of 0 as 'zero' value for being unset
		Deck: &decks.BlackjackDeck{
			DeckCount: 1,
			Deck: decks.Deck{
				Cards: []*card.Card{
					// first card out of the deck should be six
					six,
					five,
					four,
				},
			},
		},
	}

	if dt := GetCardsTotal(bj.Player.Cards); dt != 10 {
		t.Fatalf("Card total not 10, got %d", dt)
	}

	actual := bj.GetOtherMoves()

	if actual != "d" {
		t.Fatalf("Next move is wrong. Expected %s but got %s", "d", actual)
	}

}

func TestPlayerDouble(t *testing.T) {
	four, _ := card.NewCard(suit.Hearts, rank.Four, 4, true)
	six, _ := card.NewCard(suit.Hearts, rank.Six, 6, true)

	// total of 10
	playerCards := []*card.Card{
		four, six,
	}

	five, _ := card.NewCard(suit.Hearts, rank.Five, 5, true)
	seven, _ := card.NewCard(suit.Hearts, rank.Seven, 7, true)

	// total of 12
	dealerCards := []*card.Card{
		five,
		seven,
	}

	ace, _ := card.NewCard(suit.Diamonds, rank.Ace, 1, true)
	eight, _ := card.NewCard(suit.Diamonds, rank.Eight, 8, true)

	originalBet := 5
	originalCash := 500
	bj := &Blackjack{
		Player: &players.Player{
			Cards: playerCards,
			Name:  "player1",
			Cash:  originalCash,
			Bet:   originalBet,
		},
		Dealer: &players.Dealer{
			Cards: dealerCards,
		},
		// should have minCardCount of 0 as 'zero' value for being unset
		Deck: &decks.BlackjackDeck{
			DeckCount: 1,
			Deck: decks.Deck{
				Cards: []*card.Card{
					ace, // first card should be delt to player
					eight,
					five,
					four,
				},
			},
		},
	}

	actual := bj.PlayerDouble()

	if bj.Player.Bet != originalBet*2 {
		t.Fatalf("Double bet is wrong. Expected %d but got %d", originalBet, bj.Player.Bet)
	}

	if bj.Player.Cash != originalCash-originalBet {
		t.Fatalf(
			"Cash is wrong after double bet. Expected %d but got %d",
			originalCash,
			bj.Player.Cash,
		)
	}

	if actual != PlayerWon {
		t.Fatalf("Outcome is wrong. Expected %d but got %d", PlayerWon, actual)
	}

}

func TestBlackjackReturnRate(t *testing.T) {
	king, _ := card.NewCard(suit.Hearts, rank.King, 10, true)
	queen, _ := card.NewCard(suit.Hearts, rank.Queen, 10, true)

	// total 21
	playerCards := []*card.Card{
		queen,
		king,
	}

	jack, _ := card.NewCard(suit.Hearts, rank.Jack, 2, true)
	six, _ := card.NewCard(suit.Hearts, rank.Six, 6, true)

	// should be a soft 16
	dealerCards := []*card.Card{
		jack,
		six,
	}

	four, _ := card.NewCard(suit.Hearts, rank.Four, 4, true)
	five, _ := card.NewCard(suit.Hearts, rank.Five, 5, true)
	ace, _ := card.NewCard(suit.Hearts, rank.Ace, 1, true)

	bj := &Blackjack{
		Player: &players.Player{
			Cards: playerCards,
			Name:  "player1",
			Cash:  500,
			Bet:   50,
		},
		Dealer: &players.Dealer{
			Cards: dealerCards,
		},
		// should have minCardCount of 0 as 'zero' value for being unset
		Deck: &decks.BlackjackDeck{
			DeckCount: 1,
			Deck: decks.Deck{
				Cards: []*card.Card{
					ace,
					four,
					five,
				},
			},
		},
	}

	if dt := GetCardsTotal(bj.Player.Cards); dt != 20 {
		t.Fatalf("Card total wrong. Got 21, got %d", dt)
	}

	outcome := bj.PlayerHit()

	if outcome != PlayerWon {
		t.Fatalf("Round outcome is wrong. Should have %d, but got %d", PlayerWon, outcome)
	}

	if bj.payoutRate != blackjackRate {
		t.Fatalf("Payout rate is wrong. Should have %d, but got %d", blackjackRate, bj.payoutRate)
	}

	bj.PlayerWonHand()

	if bj.Player.Cash != 575 {
		t.Fatalf("Total cash is wrong after payout. Should have %d but got %d", 575, bj.Player.Cash)
	}
}
