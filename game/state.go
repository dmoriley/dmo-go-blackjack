package game

import "fmt"

type Phase string

const (
	PhaseBetting     Phase = "betting"
	PhasePlayerTurn  Phase = "player_turn"
	PhaseRoundResult Phase = "round_result"
	PhaseGameOver    Phase = "game_over"
	PhaseCashedOut   Phase = "cashed_out"
)

type Move string

const (
	MoveHit    Move = "hit"
	MoveStand  Move = "stand"
	MoveDouble Move = "double"
	MoveSplit  Move = "split"
)

type Outcome string

const (
	OutcomeWon  Outcome = "won"
	OutcomeLost Outcome = "lost"
	OutcomePush Outcome = "push"
)

type HandFinishReason string

const (
	HandFinishBust              HandFinishReason = "bust"
	HandFinishTwentyOne         HandFinishReason = "twenty_one"
	HandFinishSplitAceAutoStand HandFinishReason = "split_ace_auto_stand"
)

type CardState struct {
	Suit   string
	Rank   string
	Value  int
	FaceUp bool
}

func (c CardState) Inspect() string {
	if !c.FaceUp {
		return "{Face down}"
	}

	return fmt.Sprintf("{%s of %s, value: %d}", c.Rank, c.Suit, c.Value)
}

type HandState struct {
	Cards     []CardState
	Total     int
	Bet       int
	Done      bool
	Busted    bool
	Standing  bool
	Blackjack bool
	Resolved  bool
	Outcome   Outcome
	Payout    int
}

type Snapshot struct {
	Phase           Phase
	PlayerName      string
	Cash            int
	PreviousBet     int
	Dealer          HandState
	PlayerHands     []HandState
	ActiveHandIndex int
	IsSplitRound    bool
	LegalMoves      []Move
	DeckRemaining   int
	DeckTotal       int
	PendingResults  int
}

type StepResult struct {
	Snapshot Snapshot
	Events   []Event
}

type Event interface {
	isEvent()
}

type ReshuffledEvent struct {
	Remaining int
	Total     int
}

func (ReshuffledEvent) isEvent() {}

type HandFinishedEvent struct {
	HandIndex int
	Total     int
	Reason    HandFinishReason
}

func (HandFinishedEvent) isEvent() {}

type HandResolvedEvent struct {
	HandIndex   int
	Outcome     Outcome
	PlayerTotal int
	DealerTotal int
	Bet         int
	Payout      int
	Blackjack   bool
}

func (HandResolvedEvent) isEvent() {}

type ActionDeniedEvent struct {
	Move   Move
	Reason string
}

func (ActionDeniedEvent) isEvent() {}

type CashOutEvent struct {
	PlayerName string
	Cash       int
}

func (CashOutEvent) isEvent() {}

type GameOverEvent struct {
	PlayerName string
	Cash       int
}

func (GameOverEvent) isEvent() {}
