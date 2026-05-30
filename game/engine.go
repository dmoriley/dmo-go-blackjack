package game

import (
	"fmt"
	"slices"

	"blackjack/card"
	"blackjack/card/rank"
	"blackjack/decks"
)

const blackjackTotal = 21

type Config struct {
	PlayerName   string
	StartingCash int
	Deck         *decks.BlackjackDeck
	DeckConfig   *decks.BlackjackDeckConfig
	ShuffleCount int
}

type Engine struct {
	playerName      string
	cash            int
	previousBet     int
	deck            *decks.BlackjackDeck
	dealer          []*card.Card
	hands           []*engineHand
	activeHandIndex int
	phase           Phase
	resultQueue     []HandResolvedEvent
	resultIndex     int
}

type engineHand struct {
	cards            []*card.Card
	bet              int
	done             bool
	busted           bool
	standing         bool
	split            bool
	splitAce         bool
	naturalBlackjack bool
	resolved         bool
	outcome          Outcome
	payout           int
}

func NewEngine(config Config) (*Engine, error) {
	if config.PlayerName == "" {
		return nil, fmt.Errorf("player name is required")
	}
	if config.StartingCash < 0 {
		return nil, fmt.Errorf("starting cash cannot be negative")
	}

	startingCash := config.StartingCash
	if startingCash == 0 {
		startingCash = 500
	}

	shuffleCount := config.ShuffleCount
	if shuffleCount == 0 {
		shuffleCount = 5
	}

	deck := config.Deck
	if deck == nil {
		deckConfig := config.DeckConfig
		if deckConfig == nil {
			deckConfig = decks.NewBlackjackDeckConfig().WithNumberOfDecks(4)
		}

		deck = decks.NewBlackjackDeck(deckConfig)
		deck.Shuffle(shuffleCount)
	}

	return &Engine{
		playerName:      config.PlayerName,
		cash:            startingCash,
		deck:            deck,
		activeHandIndex: -1,
		phase:           PhaseBetting,
	}, nil
}

func (e *Engine) Snapshot() Snapshot {
	dealer := HandState{
		Cards:    snapshotCards(e.dealer),
		Total:    visibleTotal(e.dealer),
		Resolved: e.phase == PhaseRoundResult,
	}
	if e.phase == PhaseRoundResult && e.resultIndex < len(e.resultQueue) {
		dealer.Outcome = dealerOutcomeFor(e.resultQueue[e.resultIndex].Outcome)
	}

	hands := make([]HandState, 0, len(e.hands))
	for _, hand := range e.hands {
		hands = append(hands, HandState{
			Cards:     snapshotCards(hand.cards),
			Total:     fullTotal(hand.cards),
			Bet:       hand.bet,
			Done:      hand.done,
			Busted:    hand.busted,
			Standing:  hand.standing,
			Blackjack: hand.naturalBlackjack,
			Resolved:  hand.resolved,
			Outcome:   hand.outcome,
			Payout:    hand.payout,
		})
	}

	deckTotal := e.deck.DeckCount * 52
	if deckTotal == 0 {
		deckTotal = e.deck.GetLength()
	}

	pendingResults := 0
	if e.phase == PhaseRoundResult && e.resultIndex < len(e.resultQueue) {
		pendingResults = len(e.resultQueue) - e.resultIndex
	}

	return Snapshot{
		Phase:           e.phase,
		PlayerName:      e.playerName,
		Cash:            e.cash,
		PreviousBet:     e.previousBet,
		Dealer:          dealer,
		PlayerHands:     hands,
		ActiveHandIndex: e.activeSnapshotHandIndex(),
		IsSplitRound:    len(hands) > 1,
		LegalMoves:      append([]Move(nil), e.legalMoves()...),
		DeckRemaining:   e.deck.GetLength(),
		DeckTotal:       deckTotal,
		PendingResults:  pendingResults,
	}
}

func (e *Engine) PlaceBet(amount int) (StepResult, error) {
	if e.phase != PhaseBetting {
		return StepResult{}, fmt.Errorf("cannot place a bet during %s", e.phase)
	}
	if amount <= 0 {
		return StepResult{}, fmt.Errorf("bet must be greater than $0")
	}
	if amount > e.cash {
		return StepResult{}, fmt.Errorf("bet $%d exceeds available cash $%d", amount, e.cash)
	}

	e.cash -= amount
	e.previousBet = amount
	e.dealer = nil
	e.hands = []*engineHand{{bet: amount}}
	e.activeHandIndex = 0
	e.resultQueue = nil
	e.resultIndex = 0
	e.phase = PhasePlayerTurn

	e.dealToHand(e.hands[0], 1, true)
	e.dealDealer(1, true)
	e.dealToHand(e.hands[0], 1, true)
	e.dealDealer(1, false)

	events := e.collectDeckEvents(nil)

	playerNatural := e.isNaturalBlackjack(e.hands[0])
	dealerNatural := e.isDealerBlackjack()
	e.hands[0].naturalBlackjack = playerNatural

	if playerNatural || dealerNatural {
		e.revealDealer()

		var resolved HandResolvedEvent
		switch {
		case playerNatural && dealerNatural:
			resolved = e.settleHand(0, OutcomePush)
		case playerNatural:
			resolved = e.settleHand(0, OutcomeWon)
		default:
			resolved = e.settleHand(0, OutcomeLost)
		}

		return e.enterRoundResult(events, []HandResolvedEvent{resolved}), nil
	}

	return StepResult{Snapshot: e.Snapshot(), Events: events}, nil
}

func (e *Engine) ApplyMove(move Move) (StepResult, error) {
	if e.phase != PhasePlayerTurn {
		return StepResult{}, fmt.Errorf("cannot apply a move during %s", e.phase)
	}

	if !slices.Contains(e.legalMoves(), move) {
		return StepResult{}, fmt.Errorf("move %q is not legal in the current state", move)
	}

	hand := e.currentHand()
	events := []Event{}

	switch move {
	case MoveHit:
		e.dealToHand(hand, 1, true)
		events = e.collectDeckEvents(events)

		total := fullTotal(hand.cards)
		if total > blackjackTotal {
			hand.busted = true
			hand.done = true

			if len(e.hands) == 1 || e.nextPendingHand(e.activeHandIndex) == -1 {
				resolved := e.settleHand(e.activeHandIndex, OutcomeLost)
				return e.enterRoundResult(events, []HandResolvedEvent{resolved}), nil
			}

			e.markResolvedLoss(e.activeHandIndex)
			events = append(events, HandFinishedEvent{
				HandIndex: e.activeHandIndex,
				Total:     total,
				Reason:    HandFinishBust,
			})
			return e.advanceToNextHand(events)
		}

		if total == blackjackTotal {
			hand.done = true
			hand.standing = true

			if len(e.hands) == 1 || e.nextPendingHand(e.activeHandIndex) == -1 {
				return e.resolveRound(events)
			}

			events = append(events, HandFinishedEvent{
				HandIndex: e.activeHandIndex,
				Total:     total,
				Reason:    HandFinishTwentyOne,
			})
			return e.advanceToNextHand(events)
		}

		return StepResult{Snapshot: e.Snapshot(), Events: events}, nil

	case MoveStand:
		hand.done = true
		hand.standing = true
		if len(e.hands) == 1 || e.nextPendingHand(e.activeHandIndex) == -1 {
			return e.resolveRound(events)
		}

		return e.advanceToNextHand(events)

	case MoveDouble:
		if e.cash < hand.bet {
			return StepResult{
				Snapshot: e.Snapshot(),
				Events: []Event{ActionDeniedEvent{
					Move:   move,
					Reason: fmt.Sprintf("You do not have enough cash left to perform a double.\nNeed $%d, but you only have $%d", hand.bet, e.cash),
				}},
			}, nil
		}

		e.cash -= hand.bet
		hand.bet *= 2
		e.dealToHand(hand, 1, true)
		events = e.collectDeckEvents(events)
		hand.done = true
		hand.standing = true

		if fullTotal(hand.cards) > blackjackTotal {
			hand.busted = true
			resolved := e.settleHand(e.activeHandIndex, OutcomeLost)
			return e.enterRoundResult(events, []HandResolvedEvent{resolved}), nil
		}

		return e.resolveRound(events)

	case MoveSplit:
		if e.cash < hand.bet {
			return StepResult{
				Snapshot: e.Snapshot(),
				Events: []Event{ActionDeniedEvent{
					Move:   move,
					Reason: fmt.Sprintf("You dont have enough cash to split.\nNeed $%d, but you only have $%d", hand.bet, e.cash),
				}},
			}, nil
		}

		e.cash -= hand.bet

		left := &engineHand{
			cards:    []*card.Card{hand.cards[0]},
			bet:      hand.bet,
			split:    true,
			splitAce: hand.cards[0].Rank.Name == rank.Ace,
		}
		right := &engineHand{
			cards:    []*card.Card{hand.cards[1]},
			bet:      hand.bet,
			split:    true,
			splitAce: hand.cards[1].Rank.Name == rank.Ace,
		}

		e.hands = []*engineHand{left, right}
		e.activeHandIndex = 0
		return e.prepareActiveHand(events)
	}

	return StepResult{}, fmt.Errorf("unknown move %q", move)
}

func (e *Engine) Continue() (StepResult, error) {
	if e.phase != PhaseRoundResult {
		return StepResult{}, fmt.Errorf("cannot continue during %s", e.phase)
	}

	e.resultIndex++
	if e.resultIndex < len(e.resultQueue) {
		e.activeHandIndex = e.resultQueue[e.resultIndex].HandIndex
		return StepResult{
			Snapshot: e.Snapshot(),
			Events:   []Event{e.resultQueue[e.resultIndex]},
		}, nil
	}

	e.cleanupRound()
	if e.cash == 0 {
		e.phase = PhaseGameOver
		return StepResult{
			Snapshot: e.Snapshot(),
			Events: []Event{GameOverEvent{
				PlayerName: e.playerName,
				Cash:       e.cash,
			}},
		}, nil
	}

	e.phase = PhaseBetting
	return StepResult{Snapshot: e.Snapshot()}, nil
}

func (e *Engine) CashOut() (StepResult, error) {
	if e.phase != PhaseBetting {
		return StepResult{}, fmt.Errorf("can only cash out between rounds")
	}

	e.phase = PhaseCashedOut
	return StepResult{
		Snapshot: e.Snapshot(),
		Events: []Event{CashOutEvent{
			PlayerName: e.playerName,
			Cash:       e.cash,
		}},
	}, nil
}

func (e *Engine) activeSnapshotHandIndex() int {
	if len(e.hands) == 0 {
		return -1
	}

	if e.activeHandIndex < 0 || e.activeHandIndex >= len(e.hands) {
		return -1
	}

	return e.activeHandIndex
}

func (e *Engine) legalMoves() []Move {
	if e.phase != PhasePlayerTurn || e.activeHandIndex < 0 || e.activeHandIndex >= len(e.hands) {
		return nil
	}

	hand := e.currentHand()
	if hand == nil || hand.done {
		return nil
	}

	moves := []Move{MoveHit, MoveStand}
	total := fullTotal(hand.cards)

	if len(hand.cards) == 2 && !hand.split && (total == 9 || total == 10 || total == 11) {
		moves = append(moves, MoveDouble)
	}

	if len(hand.cards) == 2 && !hand.split && len(e.hands) == 1 && sameRank(hand.cards) {
		moves = append(moves, MoveSplit)
	}

	return moves
}

func (e *Engine) currentHand() *engineHand {
	if e.activeHandIndex < 0 || e.activeHandIndex >= len(e.hands) {
		return nil
	}

	return e.hands[e.activeHandIndex]
}

func (e *Engine) nextPendingHand(current int) int {
	for idx := current + 1; idx < len(e.hands); idx++ {
		if !e.hands[idx].done {
			return idx
		}
	}

	return -1
}

func (e *Engine) prepareActiveHand(events []Event) (StepResult, error) {
	for {
		hand := e.currentHand()
		if hand == nil {
			return StepResult{}, fmt.Errorf("no active hand available")
		}

		if len(hand.cards) == 1 {
			e.dealToHand(hand, 1, true)
			events = e.collectDeckEvents(events)
		}

		if hand.splitAce {
			hand.done = true
			hand.standing = true
			events = append(events, HandFinishedEvent{
				HandIndex: e.activeHandIndex,
				Total:     fullTotal(hand.cards),
				Reason:    HandFinishSplitAceAutoStand,
			})
			next := e.nextPendingHand(e.activeHandIndex)
			if next == -1 {
				return e.resolveRound(events)
			}

			e.activeHandIndex = next
			continue
		}

		e.phase = PhasePlayerTurn
		return StepResult{Snapshot: e.Snapshot(), Events: events}, nil
	}
}

func (e *Engine) advanceToNextHand(events []Event) (StepResult, error) {
	next := e.nextPendingHand(e.activeHandIndex)
	if next == -1 {
		return e.resolveRound(events)
	}

	e.activeHandIndex = next
	return e.prepareActiveHand(events)
}

func (e *Engine) resolveRound(events []Event) (StepResult, error) {
	if e.hasUnresolvedLiveHand() {
		e.revealDealer()
		for e.dealerShouldHit() {
			e.dealDealer(1, true)
			events = e.collectDeckEvents(events)
		}
	}

	resolved := []HandResolvedEvent{}
	dealerTotal := fullTotal(e.dealer)
	for idx, hand := range e.hands {
		if hand.resolved {
			continue
		}

		playerTotal := fullTotal(hand.cards)
		outcome := OutcomePush
		switch {
		case hand.busted:
			outcome = OutcomeLost
		case dealerTotal > blackjackTotal:
			outcome = OutcomeWon
		case playerTotal > dealerTotal:
			outcome = OutcomeWon
		case playerTotal < dealerTotal:
			outcome = OutcomeLost
		default:
			outcome = OutcomePush
		}

		resolved = append(resolved, e.settleHand(idx, outcome))
	}

	return e.enterRoundResult(events, resolved), nil
}

func (e *Engine) hasUnresolvedLiveHand() bool {
	for _, hand := range e.hands {
		if !hand.resolved && !hand.busted {
			return true
		}
	}

	return false
}

func (e *Engine) settleHand(index int, outcome Outcome) HandResolvedEvent {
	hand := e.hands[index]
	hand.done = true
	hand.resolved = true
	hand.outcome = outcome

	payout := 0
	switch outcome {
	case OutcomeWon:
		if hand.naturalBlackjack {
			payout = hand.bet * 25 / 10
		} else {
			payout = hand.bet * 2
		}
		e.cash += payout
	case OutcomePush:
		payout = hand.bet
		e.cash += payout
	}

	hand.payout = payout

	return HandResolvedEvent{
		HandIndex:   index,
		Outcome:     outcome,
		PlayerTotal: fullTotal(hand.cards),
		DealerTotal: fullTotal(e.dealer),
		Bet:         hand.bet,
		Payout:      payout,
		Blackjack:   hand.naturalBlackjack,
	}
}

func (e *Engine) markResolvedLoss(index int) {
	hand := e.hands[index]
	hand.done = true
	hand.resolved = true
	hand.outcome = OutcomeLost
	hand.payout = 0
	if hand.busted {
		hand.standing = false
	}
}

func (e *Engine) enterRoundResult(events []Event, resolved []HandResolvedEvent) StepResult {
	e.resultQueue = resolved
	e.resultIndex = 0
	e.phase = PhaseRoundResult
	e.activeHandIndex = resolved[0].HandIndex

	resultEvents := make([]Event, 0, len(events)+1)
	resultEvents = append(resultEvents, events...)
	resultEvents = append(resultEvents, resolved[0])

	return StepResult{
		Snapshot: e.Snapshot(),
		Events:   resultEvents,
	}
}

func (e *Engine) cleanupRound() {
	for _, hand := range e.hands {
		e.deck.AddDiscardedCards(hand.cards)
	}
	e.deck.AddDiscardedCards(e.dealer)

	e.dealer = nil
	e.hands = nil
	e.activeHandIndex = -1
	e.resultQueue = nil
	e.resultIndex = 0
}

func dealerOutcomeFor(playerOutcome Outcome) Outcome {
	switch playerOutcome {
	case OutcomeWon:
		return OutcomeLost
	case OutcomeLost:
		return OutcomeWon
	default:
		return OutcomePush
	}
}

func (e *Engine) revealDealer() {
	for _, dealtCard := range e.dealer {
		dealtCard.IsFaceUp = true
	}
}

func (e *Engine) dealerShouldHit() bool {
	total, soft := totalAndSoft(e.dealer, true)
	return total < 17 || (total == 17 && soft)
}

func (e *Engine) isDealerBlackjack() bool {
	return len(e.dealer) == 2 && fullTotal(e.dealer) == blackjackTotal
}

func (e *Engine) isNaturalBlackjack(hand *engineHand) bool {
	return !hand.split && len(hand.cards) == 2 && fullTotal(hand.cards) == blackjackTotal
}

func (e *Engine) dealToHand(hand *engineHand, count int, faceUp bool) {
	for _, dealtCard := range e.deck.Pop(count) {
		dealtCard.IsFaceUp = faceUp
		hand.cards = append(hand.cards, dealtCard)
	}
}

func (e *Engine) dealDealer(count int, faceUp bool) {
	for _, dealtCard := range e.deck.Pop(count) {
		dealtCard.IsFaceUp = faceUp
		e.dealer = append(e.dealer, dealtCard)
	}
}

func (e *Engine) collectDeckEvents(events []Event) []Event {
	if e.deck.ConsumeReshuffle() {
		return append(events, ReshuffledEvent{
			Remaining: e.deck.GetLength(),
			Total:     e.deck.DeckCount * 52,
		})
	}

	return events
}

func snapshotCards(cards []*card.Card) []CardState {
	states := make([]CardState, 0, len(cards))
	for _, dealtCard := range cards {
		states = append(states, CardState{
			Suit:   dealtCard.Suit,
			Rank:   dealtCard.Rank.Name,
			Value:  dealtCard.Rank.Value,
			FaceUp: dealtCard.IsFaceUp,
		})
	}

	return states
}

func visibleTotal(cards []*card.Card) int {
	total, _ := totalAndSoft(cards, false)
	return total
}

func fullTotal(cards []*card.Card) int {
	total, _ := totalAndSoft(cards, true)
	return total
}

func totalAndSoft(cards []*card.Card, includeFaceDown bool) (int, bool) {
	nonAceTotal := 0
	aceCount := 0

	for _, dealtCard := range cards {
		if !includeFaceDown && !dealtCard.IsFaceUp {
			continue
		}

		if dealtCard.Rank.Name == rank.Ace {
			aceCount++
			continue
		}

		nonAceTotal += dealtCard.Rank.Value
	}

	total := nonAceTotal + aceCount
	soft := false
	if aceCount > 0 && total+10 <= blackjackTotal {
		total += 10
		soft = true
	}

	return total, soft
}

func sameRank(cards []*card.Card) bool {
	return len(cards) == 2 && cards[0].Rank.Name == cards[1].Rank.Name
}
