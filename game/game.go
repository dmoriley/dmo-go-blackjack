// Package game contains engine logic to run the blackjack game
package game

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"blackjack/decks"
	"blackjack/game/utils"
)

const PROMPT = "=> "

const (
	HIT    = "h"
	STAND  = "s"
	DOUBLE = "d"
	SPLIT  = "l"
)

func Start(in io.Reader, out io.Writer, playerName string) error {
	config := decks.NewBlackjackDeckConfig().WithNumberOfDecks(4)
	engine, err := NewEngine(Config{
		PlayerName:   playerName,
		StartingCash: 500,
		DeckConfig:   config,
		ShuffleCount: 5,
	})
	if err != nil {
		return err
	}

	return startWithEngine(in, out, engine)
}

func startWithEngine(in io.Reader, out io.Writer, engine *Engine) error {
	scanner := bufio.NewScanner(in)
	var pending *StepResult

	for {
		snapshot := engine.Snapshot()
		if pending != nil {
			snapshot = pending.Snapshot
		}

		switch snapshot.Phase {
		case PhaseBetting:
			if pending != nil {
				printEvents(out, pending.Events)
				pending = nil
			}

			if snapshot.Cash == 0 {
				fmt.Fprintln(out, "\n*******************************")
				fmt.Fprintln(out, "Busted! You're out of money.")
				fmt.Fprintln(out, "\n*******************************")
				fmt.Fprintln(out, "Thanks for playing! Come back with more cash.")
				return nil
			}

			fmt.Fprintf(out, "%s has $%d in wallet\n", snapshot.PlayerName, snapshot.Cash)
			for {
				fmt.Fprint(out, "Place your bet or \".cashout\": $")
				text, err := readInput(scanner)
				if err != nil {
					fmt.Fprintln(out, err.Error())
					fmt.Fprintln(out, "That bet is invalid. Please try again.")
					continue
				}

				if strings.HasPrefix(text, ".") {
					result, shouldContinue, err := handleBetCommand(out, engine, snapshot.PreviousBet, text)
					if err != nil {
						fmt.Fprintln(out, err.Error())
						continue
					}
					if !shouldContinue {
						pending = result
						break
					}
					if text == ".help" || text == ".h" {
						continue
					}
					if text == "." || text == ".." {
						text = strconv.Itoa(snapshot.PreviousBet)
					}
				}

				bet, err := strconv.Atoi(text)
				if err != nil {
					fmt.Fprintln(out, err.Error())
					fmt.Fprintln(out, "That bet is invalid. Please try again.")
					continue
				}

				result, err := engine.PlaceBet(bet)
				if err != nil {
					fmt.Fprintln(out, err.Error())
					continue
				}

				pending = &result
				break
			}

		case PhasePlayerTurn:
			renderSnapshot(out, snapshot, tableTitle(snapshot), activeSubtitle(snapshot))
			if pending != nil {
				printEvents(out, pending.Events)
				pending = nil
			}
			printMovePrompt(out, snapshot.LegalMoves)
			for {
				fmt.Fprint(out, PROMPT)
				text, err := readInput(scanner)
				if err != nil {
					fmt.Fprintln(out, err.Error())
					continue
				}

				move, ok := parseMove(text)
				if !ok {
					fmt.Fprintf(out, "Unexpected value entered: %q\n", text)
					continue
				}

				result, err := engine.ApplyMove(move)
				if err != nil {
					fmt.Fprintln(out, err.Error())
					continue
				}

				pending = &result
				break
			}

		case PhaseRoundResult:
			renderSnapshot(out, snapshot, tableTitle(snapshot), activeSubtitle(snapshot))
			if pending != nil {
				printEvents(out, pending.Events)
				pending = nil
			}
			fmt.Fprint(out, "\nPress ENTER to continue...")
			if _, err := utils.GetUserInput(utils.NewInputConfig(scanner).SetAnyKey()); err != nil {
				return err
			}

			result, err := engine.Continue()
			if err != nil {
				return err
			}
			pending = &result

		case PhaseGameOver:
			if pending != nil {
				printEvents(out, pending.Events)
				pending = nil
			}
			fmt.Fprintln(out, "\n*******************************")
			fmt.Fprintln(out, "Busted! You're out of money.")
			fmt.Fprintln(out, "\n*******************************")
			fmt.Fprintln(out, "Thanks for playing! Come back with more cash.")
			return nil

		case PhaseCashedOut:
			if pending != nil {
				printEvents(out, pending.Events)
				pending = nil
			}
			fmt.Fprintf(out, "\nCongratulations %s! You're cashing out with $%d\n", snapshot.PlayerName, snapshot.Cash)
			fmt.Fprintln(out, "\n*******************************")
			fmt.Fprintln(out, "Thanks for playing! Come back with more cash.")
			return nil
		}
	}
}

func tableTitle(snapshot Snapshot) string {
	if snapshot.IsSplitRound {
		return "Table Split Round"
	}
	return "Table Cards"
}

func activeSubtitle(snapshot Snapshot) string {
	if snapshot.IsSplitRound && snapshot.ActiveHandIndex >= 0 {
		return fmt.Sprintf("Hand: %d", snapshot.ActiveHandIndex+1)
	}
	return ""
}

func handleBetCommand(out io.Writer, engine *Engine, previousBet int, command string) (*StepResult, bool, error) {
	switch command {
	case ".help", ".h":
		fmt.Fprintln(out, `
	-------------------------------------
	| LAST BET |    CASHOUT   |   HELP  |
	|    (.)   |  (.cashout)  | (.help) |
	-------------------------------------
		`)
		return nil, true, nil
	case ".cashout", ".c", ".exit", ".quit":
		result, err := engine.CashOut()
		return &result, false, err
	case ".", "..":
		if previousBet == 0 {
			return nil, true, fmt.Errorf("no previous bet available")
		}
		return nil, true, nil
	default:
		return nil, true, fmt.Errorf("Invalid command. Try \".help\" for more options")
	}
}

func readInput(scanner *bufio.Scanner) (string, error) {
	config := utils.NewInputConfig(scanner)
	return utils.GetUserInput(config)
}

func parseMove(input string) (Move, bool) {
	switch input {
	case HIT:
		return MoveHit, true
	case STAND:
		return MoveStand, true
	case DOUBLE:
		return MoveDouble, true
	case SPLIT:
		return MoveSplit, true
	default:
		return "", false
	}
}

func printMovePrompt(out io.Writer, moves []Move) {
	hasDouble := slicesContainsMove(moves, MoveDouble)
	hasSplit := slicesContainsMove(moves, MoveSplit)

	switch {
	case hasDouble && hasSplit:
		fmt.Fprintln(out, `
	--------------------------------
	| HIT | STAND | DOUBLE | SPLIT |
	| (h) |  (s)  |  (d)   |  (l)  |
	--------------------------------
		`)
	case hasDouble:
		fmt.Fprintln(out, `
	------------------------
	| HIT | STAND | DOUBLE |
	| (h) |  (s)  |  (d)   |
	------------------------
		`)
	case hasSplit:
		fmt.Fprintln(out, `
	-----------------------
	| HIT | STAND | SPLIT |
	| (h) |  (s)  |  (l)  |
	-----------------------
		`)
	default:
		fmt.Fprintln(out, `
	---------------
	| HIT | STAND |
	| (h) |  (s)  |
	---------------
		`)
	}
}

func printEvents(out io.Writer, events []Event) {
	for _, event := range events {
		switch e := event.(type) {
		case ReshuffledEvent:
			fmt.Fprintln(out, "\n***********************************")
			fmt.Fprintln(out, "***** Reshuffling the deck... *****")
			fmt.Fprintln(out, "***********************************")
		case HandFinishedEvent:
			switch e.Reason {
			case HandFinishBust:
				fmt.Fprint(out, "***  Dealer win!  ***\nCollecting all losing bets...\n\n")
			case HandFinishTwentyOne:
				fmt.Fprintln(out, "Player hit max card score.")
			case HandFinishSplitAceAutoStand:
				fmt.Fprintln(out, "Split aces receive one card and stand automatically.")
			}
		case HandResolvedEvent:
			switch e.Outcome {
			case OutcomeWon:
				if e.Blackjack {
					fmt.Fprintln(out, "BLACKJACK!!\nCollect your winnings at a rate of 1.5.")
				}
				fmt.Fprint(out, "***  Player win!  ***\nAdding winnings to your wallet...\n\n")
			case OutcomeLost:
				fmt.Fprint(out, "***  Dealer win!  ***\nCollecting all losing bets...\n\n")
			case OutcomePush:
				fmt.Fprint(out, "Push! Returning all bets...\n\n")
			}
		case ActionDeniedEvent:
			fmt.Fprintln(out, e.Reason)
		}
	}
}

func slicesContainsMove(moves []Move, target Move) bool {
	for _, move := range moves {
		if move == target {
			return true
		}
	}
	return false
}

func renderSnapshot(out io.Writer, snapshot Snapshot, title string, subtitle string) {
	playerHand := HandState{}
	if snapshot.ActiveHandIndex >= 0 && snapshot.ActiveHandIndex < len(snapshot.PlayerHands) {
		playerHand = snapshot.PlayerHands[snapshot.ActiveHandIndex]
	} else if len(snapshot.PlayerHands) > 0 {
		playerHand = snapshot.PlayerHands[0]
	}

	config := utils.NewPrintTableConfig(
		toTableHand("Dealer", snapshot.Dealer),
		toTableHand(snapshot.PlayerName, playerHand),
		snapshot.DeckRemaining,
		snapshot.DeckTotal,
	)
	if title != "" {
		config.SetTitle(title)
	}
	if subtitle != "" {
		config.SetSubtitle(subtitle)
	}

	fmt.Fprint(out, utils.RenderTable(config))
}

func toTableHand(label string, hand HandState) utils.TableHand {
	cards := make([]utils.TableCard, 0, len(hand.Cards))
	for _, card := range hand.Cards {
		cards = append(cards, utils.TableCard{
			Label:  fmt.Sprintf("%s of %s", card.Rank, card.Suit),
			Value:  card.Value,
			FaceUp: card.FaceUp,
		})
	}

	return utils.TableHand{
		Label: label,
		Total: hand.Total,
		Cards: cards,
	}
}
