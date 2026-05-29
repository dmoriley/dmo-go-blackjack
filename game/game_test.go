package game

import (
	"bytes"
	"strings"
	"testing"

	"blackjack/card"
	"blackjack/card/rank"
	"blackjack/card/suit"
	"blackjack/decks"
)

func TestNaturalBlackjackPaysThreeToTwo(t *testing.T) {
	engine := newTestEngine(t, 500,
		newCard(t, suit.Hearts, rank.Ace, 1, false),
		newCard(t, suit.Spades, rank.Five, 5, false),
		newCard(t, suit.Diamonds, rank.King, 10, false),
		newCard(t, suit.Clubs, rank.Nine, 9, false),
	)

	result, err := engine.PlaceBet(50)
	if err != nil {
		t.Fatalf("PlaceBet returned error: %v", err)
	}

	assertSnapshotPhase(t, result.Snapshot, PhaseRoundResult)
	assertHandResolvedEvent(t, result.Events, HandResolvedEvent{
		HandIndex:   0,
		Outcome:     OutcomeWon,
		PlayerTotal: 21,
		DealerTotal: 14,
		Bet:         50,
		Payout:      125,
		Blackjack:   true,
	})

	if result.Snapshot.Cash != 575 {
		t.Fatalf("cash wrong after natural blackjack. want 575 got %d", result.Snapshot.Cash)
	}
}

func TestHitToTwentyOnePaysEvenMoney(t *testing.T) {
	engine := newTestEngine(t, 500,
		newCard(t, suit.Hearts, rank.Queen, 10, false),
		newCard(t, suit.Spades, rank.Six, 6, false),
		newCard(t, suit.Diamonds, rank.King, 10, false),
		newCard(t, suit.Clubs, rank.Five, 5, false),
		newCard(t, suit.Hearts, rank.Ace, 1, false),
		newCard(t, suit.Spades, rank.Six, 6, false),
	)

	if _, err := engine.PlaceBet(50); err != nil {
		t.Fatalf("PlaceBet returned error: %v", err)
	}

	result, err := engine.ApplyMove(MoveHit)
	if err != nil {
		t.Fatalf("ApplyMove returned error: %v", err)
	}

	assertSnapshotPhase(t, result.Snapshot, PhaseRoundResult)
	resolved := assertResolvedEventOutcome(t, result.Events, OutcomeWon)
	if resolved.Blackjack {
		t.Fatalf("hit to 21 should not count as blackjack")
	}
	if resolved.Payout != 100 {
		t.Fatalf("hit to 21 payout wrong. want 100 got %d", resolved.Payout)
	}
	if result.Snapshot.Cash != 550 {
		t.Fatalf("cash wrong after hit to 21. want 550 got %d", result.Snapshot.Cash)
	}
}

func TestDealerBlackjackBeatsPlayer(t *testing.T) {
	engine := newTestEngine(t, 500,
		newCard(t, suit.Hearts, rank.Ten, 10, false),
		newCard(t, suit.Spades, rank.Ace, 1, false),
		newCard(t, suit.Diamonds, rank.Nine, 9, false),
		newCard(t, suit.Clubs, rank.King, 10, false),
	)

	result, err := engine.PlaceBet(50)
	if err != nil {
		t.Fatalf("PlaceBet returned error: %v", err)
	}

	resolved := assertResolvedEventOutcome(t, result.Events, OutcomeLost)
	if resolved.Blackjack {
		t.Fatalf("dealer blackjack should not mark player blackjack")
	}
	if result.Snapshot.Cash != 450 {
		t.Fatalf("cash wrong after dealer blackjack. want 450 got %d", result.Snapshot.Cash)
	}
}

func TestDealerBlackjackPushesPlayerNatural(t *testing.T) {
	engine := newTestEngine(t, 500,
		newCard(t, suit.Hearts, rank.Ace, 1, false),
		newCard(t, suit.Spades, rank.Ace, 1, false),
		newCard(t, suit.Diamonds, rank.King, 10, false),
		newCard(t, suit.Clubs, rank.King, 10, false),
	)

	result, err := engine.PlaceBet(50)
	if err != nil {
		t.Fatalf("PlaceBet returned error: %v", err)
	}

	resolved := assertResolvedEventOutcome(t, result.Events, OutcomePush)
	if resolved.Payout != 50 {
		t.Fatalf("push payout wrong. want 50 got %d", resolved.Payout)
	}
	if result.Snapshot.Cash != 500 {
		t.Fatalf("cash wrong after blackjack push. want 500 got %d", result.Snapshot.Cash)
	}
}

func TestDealerHitsSoftSeventeen(t *testing.T) {
	engine := newTestEngine(t, 500,
		newCard(t, suit.Hearts, rank.Jack, 10, false),
		newCard(t, suit.Spades, rank.Ace, 1, false),
		newCard(t, suit.Diamonds, rank.King, 10, false),
		newCard(t, suit.Clubs, rank.Six, 6, false),
		newCard(t, suit.Hearts, rank.Five, 5, false),
		newCard(t, suit.Spades, rank.Nine, 9, false),
	)

	if _, err := engine.PlaceBet(5); err != nil {
		t.Fatalf("PlaceBet returned error: %v", err)
	}

	result, err := engine.ApplyMove(MoveStand)
	if err != nil {
		t.Fatalf("ApplyMove returned error: %v", err)
	}

	resolved := assertResolvedEventOutcome(t, result.Events, OutcomeLost)
	if resolved.DealerTotal != 21 {
		t.Fatalf("dealer total wrong after soft 17 draw sequence. want 21 got %d", resolved.DealerTotal)
	}
	if len(result.Snapshot.Dealer.Cards) != 4 {
		t.Fatalf("dealer should have drawn to 4 cards, got %d", len(result.Snapshot.Dealer.Cards))
	}
	if result.Snapshot.Dealer.Total != 21 {
		t.Fatalf("dealer snapshot total wrong. want 21 got %d", result.Snapshot.Dealer.Total)
	}
}

func TestDoubleMoveAvailableOnNineToElevenOnly(t *testing.T) {
	engine := newTestEngine(t, 500,
		newCard(t, suit.Hearts, rank.Four, 4, false),
		newCard(t, suit.Spades, rank.Five, 5, false),
		newCard(t, suit.Diamonds, rank.Six, 6, false),
		newCard(t, suit.Clubs, rank.Seven, 7, false),
	)

	result, err := engine.PlaceBet(5)
	if err != nil {
		t.Fatalf("PlaceBet returned error: %v", err)
	}

	assertLegalMoves(t, result.Snapshot.LegalMoves, MoveHit, MoveStand, MoveDouble)
}

func TestPlayerDoubleWinsAndDoublesBet(t *testing.T) {
	engine := newTestEngine(t, 500,
		newCard(t, suit.Hearts, rank.Four, 4, false),
		newCard(t, suit.Spades, rank.Five, 5, false),
		newCard(t, suit.Diamonds, rank.Six, 6, false),
		newCard(t, suit.Clubs, rank.Seven, 7, false),
		newCard(t, suit.Hearts, rank.Ace, 1, false),
		newCard(t, suit.Spades, rank.Five, 5, false),
	)

	if _, err := engine.PlaceBet(5); err != nil {
		t.Fatalf("PlaceBet returned error: %v", err)
	}

	result, err := engine.ApplyMove(MoveDouble)
	if err != nil {
		t.Fatalf("ApplyMove returned error: %v", err)
	}

	resolved := assertResolvedEventOutcome(t, result.Events, OutcomeWon)
	if resolved.Bet != 10 {
		t.Fatalf("double bet wrong. want 10 got %d", resolved.Bet)
	}
	if result.Snapshot.Cash != 510 {
		t.Fatalf("cash wrong after winning double. want 510 got %d", result.Snapshot.Cash)
	}
}

func TestDoubleDeniedWithoutEnoughCash(t *testing.T) {
	engine := newTestEngine(t, 5,
		newCard(t, suit.Hearts, rank.Four, 4, false),
		newCard(t, suit.Spades, rank.Five, 5, false),
		newCard(t, suit.Diamonds, rank.Six, 6, false),
		newCard(t, suit.Clubs, rank.Seven, 7, false),
	)

	if _, err := engine.PlaceBet(5); err != nil {
		t.Fatalf("PlaceBet returned error: %v", err)
	}

	result, err := engine.ApplyMove(MoveDouble)
	if err != nil {
		t.Fatalf("ApplyMove returned error: %v", err)
	}

	assertSnapshotPhase(t, result.Snapshot, PhasePlayerTurn)
	if len(result.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(result.Events))
	}
	denied, ok := result.Events[0].(ActionDeniedEvent)
	if !ok {
		t.Fatalf("expected ActionDeniedEvent, got %T", result.Events[0])
	}
	if denied.Move != MoveDouble {
		t.Fatalf("wrong denied move. want %q got %q", MoveDouble, denied.Move)
	}
}

func TestSplitMoveAvailable(t *testing.T) {
	engine := newTestEngine(t, 500,
		newCard(t, suit.Hearts, rank.Four, 4, false),
		newCard(t, suit.Spades, rank.Five, 5, false),
		newCard(t, suit.Diamonds, rank.Four, 4, false),
		newCard(t, suit.Clubs, rank.Seven, 7, false),
	)

	result, err := engine.PlaceBet(5)
	if err != nil {
		t.Fatalf("PlaceBet returned error: %v", err)
	}

	assertLegalMoves(t, result.Snapshot.LegalMoves, MoveHit, MoveStand, MoveSplit)
}

func TestSplitAndDoubleMoveAvailableTogether(t *testing.T) {
	engine := newTestEngine(t, 500,
		newCard(t, suit.Hearts, rank.Five, 5, false),
		newCard(t, suit.Spades, rank.Five, 5, false),
		newCard(t, suit.Diamonds, rank.Five, 5, false),
		newCard(t, suit.Clubs, rank.Seven, 7, false),
	)

	result, err := engine.PlaceBet(5)
	if err != nil {
		t.Fatalf("PlaceBet returned error: %v", err)
	}

	assertLegalMoves(t, result.Snapshot.LegalMoves, MoveHit, MoveStand, MoveDouble, MoveSplit)
}

func TestSplitRoundDoesNotAllowDouble(t *testing.T) {
	engine := newTestEngine(t, 500,
		newCard(t, suit.Hearts, rank.Eight, 8, false),
		newCard(t, suit.Spades, rank.Six, 6, false),
		newCard(t, suit.Diamonds, rank.Eight, 8, false),
		newCard(t, suit.Clubs, rank.Nine, 9, false),
		newCard(t, suit.Hearts, rank.Three, 3, false),
		newCard(t, suit.Spades, rank.Four, 4, false),
	)

	if _, err := engine.PlaceBet(10); err != nil {
		t.Fatalf("PlaceBet returned error: %v", err)
	}

	result, err := engine.ApplyMove(MoveSplit)
	if err != nil {
		t.Fatalf("ApplyMove returned error: %v", err)
	}

	assertSnapshotPhase(t, result.Snapshot, PhasePlayerTurn)
	if result.Snapshot.ActiveHandIndex != 0 {
		t.Fatalf("expected first split hand active, got %d", result.Snapshot.ActiveHandIndex)
	}
	assertLegalMoves(t, result.Snapshot.LegalMoves, MoveHit, MoveStand)
}

func TestSplitAcesAutoStandAfterOneCardEach(t *testing.T) {
	engine := newTestEngine(t, 500,
		newCard(t, suit.Hearts, rank.Ace, 1, false),
		newCard(t, suit.Spades, rank.Six, 6, false),
		newCard(t, suit.Diamonds, rank.Ace, 1, false),
		newCard(t, suit.Clubs, rank.Eight, 8, false),
		newCard(t, suit.Hearts, rank.Nine, 9, false),
		newCard(t, suit.Spades, rank.Two, 2, false),
		newCard(t, suit.Diamonds, rank.Seven, 7, false),
	)

	if _, err := engine.PlaceBet(10); err != nil {
		t.Fatalf("PlaceBet returned error: %v", err)
	}

	result, err := engine.ApplyMove(MoveSplit)
	if err != nil {
		t.Fatalf("ApplyMove returned error: %v", err)
	}

	assertSnapshotPhase(t, result.Snapshot, PhaseRoundResult)
	if len(result.Snapshot.PlayerHands) != 2 {
		t.Fatalf("expected 2 split hands, got %d", len(result.Snapshot.PlayerHands))
	}
	for idx, hand := range result.Snapshot.PlayerHands {
		if len(hand.Cards) != 2 {
			t.Fatalf("split ace hand %d should have 2 cards, got %d", idx, len(hand.Cards))
		}
	}
	assertEventReason(t, result.Events, HandFinishSplitAceAutoStand)
}

func TestReshuffleSurfacesAsEvent(t *testing.T) {
	deck := decks.NewBlackjackDeck(decks.NewBlackjackDeckConfig().WithNumberOfDecks(1).WithMinCardCount(2))
	used := deck.Cards[:5]
	deck.Cards = deck.Cards[5:]
	for _, dealtCard := range used {
		dealtCard.IsFaceUp = true
	}
	deck.AddDiscardedCards(used)
	deck.Cards = deck.Cards[:6]

	engine, err := NewEngine(Config{
		PlayerName:   "Tester",
		StartingCash: 500,
		Deck:         deck,
	})
	if err != nil {
		t.Fatalf("NewEngine returned error: %v", err)
	}

	result, err := engine.PlaceBet(5)
	if err != nil {
		t.Fatalf("PlaceBet returned error: %v", err)
	}

	assertHasReshuffleEvent(t, result.Events)
}

func TestGameOverAfterLosingLastCash(t *testing.T) {
	engine := newTestEngine(t, 5,
		newCard(t, suit.Hearts, rank.Ten, 10, false),
		newCard(t, suit.Spades, rank.Ten, 10, false),
		newCard(t, suit.Diamonds, rank.Six, 6, false),
		newCard(t, suit.Clubs, rank.King, 10, false),
	)

	if _, err := engine.PlaceBet(5); err != nil {
		t.Fatalf("PlaceBet returned error: %v", err)
	}

	result, err := engine.ApplyMove(MoveStand)
	if err != nil {
		t.Fatalf("ApplyMove returned error: %v", err)
	}
	assertResolvedEventOutcome(t, result.Events, OutcomeLost)

	continueResult, err := engine.Continue()
	if err != nil {
		t.Fatalf("Continue returned error: %v", err)
	}

	assertSnapshotPhase(t, continueResult.Snapshot, PhaseGameOver)
	if len(continueResult.Events) != 1 {
		t.Fatalf("expected one game-over event, got %d", len(continueResult.Events))
	}
	if _, ok := continueResult.Events[0].(GameOverEvent); !ok {
		t.Fatalf("expected GameOverEvent, got %T", continueResult.Events[0])
	}
}

func TestStartRendersRoundStatesOnce(t *testing.T) {
	engine := newTestEngine(t, 500,
		newCard(t, suit.Hearts, rank.Ten, 10, false),
		newCard(t, suit.Spades, rank.Nine, 9, false),
		newCard(t, suit.Diamonds, rank.Eight, 8, false),
		newCard(t, suit.Clubs, rank.Seven, 7, false),
		newCard(t, suit.Hearts, rank.King, 10, false),
	)

	input := strings.NewReader("5\ns\n\n.cashout\n")
	var output bytes.Buffer

	if err := startWithEngine(input, &output, engine); err != nil {
		t.Fatalf("startWithEngine returned error: %v", err)
	}

	if got := strings.Count(output.String(), "Table Cards"); got != 2 {
		t.Fatalf("expected 2 table renders for one round, got %d\noutput:\n%s", got, output.String())
	}
	if got := strings.Count(output.String(), "| HIT | STAND |"); got != 1 {
		t.Fatalf("expected move prompt once, got %d\noutput:\n%s", got, output.String())
	}
}

func newTestEngine(t *testing.T, startingCash int, cards ...*card.Card) *Engine {
	t.Helper()

	deck := &decks.BlackjackDeck{
		DeckCount: 1,
		Deck: decks.Deck{
			Cards: cards,
		},
	}

	engine, err := NewEngine(Config{
		PlayerName:   "Tester",
		StartingCash: startingCash,
		Deck:         deck,
	})
	if err != nil {
		t.Fatalf("NewEngine returned error: %v", err)
	}

	return engine
}

func newCard(t *testing.T, suitName string, rankName string, value int, faceUp bool) *card.Card {
	t.Helper()
	dealtCard, err := card.NewCard(suitName, rankName, value, faceUp)
	if err != nil {
		t.Fatalf("NewCard returned error: %v", err)
	}
	return dealtCard
}

func assertSnapshotPhase(t *testing.T, snapshot Snapshot, want Phase) {
	t.Helper()
	if snapshot.Phase != want {
		t.Fatalf("wrong phase. want %q got %q", want, snapshot.Phase)
	}
}

func assertLegalMoves(t *testing.T, got []Move, want ...Move) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("wrong legal moves length. want %d got %d (%v)", len(want), len(got), got)
	}
	for idx := range want {
		if got[idx] != want[idx] {
			t.Fatalf("wrong legal move at %d. want %q got %q", idx, want[idx], got[idx])
		}
	}
}

func assertEventReason(t *testing.T, events []Event, want HandFinishReason) {
	t.Helper()
	for _, event := range events {
		finished, ok := event.(HandFinishedEvent)
		if ok && finished.Reason == want {
			return
		}
	}
	t.Fatalf("expected HandFinishedEvent with reason %q", want)
}

func assertResolvedEventOutcome(t *testing.T, events []Event, want Outcome) HandResolvedEvent {
	t.Helper()
	for _, event := range events {
		resolved, ok := event.(HandResolvedEvent)
		if ok && resolved.Outcome == want {
			return resolved
		}
	}
	t.Fatalf("expected HandResolvedEvent with outcome %q", want)
	return HandResolvedEvent{}
}

func assertHandResolvedEvent(t *testing.T, events []Event, want HandResolvedEvent) {
	t.Helper()
	resolved := assertResolvedEventOutcome(t, events, want.Outcome)
	if resolved != want {
		t.Fatalf("wrong resolved event. want %+v got %+v", want, resolved)
	}
}

func assertHasReshuffleEvent(t *testing.T, events []Event) {
	t.Helper()
	for _, event := range events {
		if _, ok := event.(ReshuffledEvent); ok {
			return
		}
	}
	t.Fatalf("expected reshuffle event")
}
