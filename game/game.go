package game

import (
	"blackjack/card"
	"blackjack/card/rank"
	"blackjack/decks"
	"blackjack/game/players"
	"blackjack/game/utils"
	"bufio"
	"fmt"
	"io"
	"strings"
)

const (
	// Max value to win or bust
	BLACKJACK = 21
	PROMPT    = "=> "
)

// Start the game
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	player := &players.Player{
		// TODO: change player name to player input value at start of game
		Name: "David",
		Cash: 500,
	}
	dealer := &players.Dealer{}
	// TODO: change deck count to a player input at config of game
	config := decks.NewBlackjackDeckConfig().WithNumberOfDecks(4)
	deck := decks.NewBlackjackDeck(config)
	deck.Shuffle(5)

	blackjack := &Blackjack{
		Dealer:     dealer,
		Player:     player,
		Deck:       deck,
		scanner:    scanner,
		payoutRate: normalRate,
	}

	for {
		outcome := InProgress
		blackjack.PlaceBet(blackjack.Player)
		blackjack.DealFirstCards()
		if playerScore := utils.CalcCardsTotal(blackjack.Player.Cards); playerScore >= 21 {
			// player hit 21 points or is over, in either case
			// they have no more moves to play, so stand
			if playerScore == BLACKJACK {
				blackjack.payoutRate = blackjackRate
				outcome = blackjack.DealerBlackjackCheck()
			} else {
				outcome = blackjack.PlayerStand()
			}
		}

		for outcome == InProgress {
			for outcome == InProgress {

				move := blackjack.ChooseNextMove()

				switch move {
				case HIT:
					outcome = blackjack.PlayerHit()
				case STAND:
					outcome = blackjack.PlayerStand()
				case DOUBLE:
					outcome = blackjack.PlayerDouble()
				case SPLIT:
					outcome = blackjack.PlayerSplit()
				}

			}

			switch outcome {
			case PlayerWon:
				blackjack.PlayerWonHand()
			case PlayerLost:
				blackjack.PlayerLostHand()
			case Standoff:
				blackjack.Standoff()
			}
			// case done skips the switch

			blackjack.cleanup()
		}

		continuePrompt := `
	----------------------
	| CASHOUT | CONTINUE |  
	|   (c)   |  (enter) |
	----------------------
		`
		fmt.Println(continuePrompt)
		fmt.Print(PROMPT)
		continueConfig := utils.NewInputConfig(blackjack.scanner).SetAnyKey()
		res, _ := utils.GetUserInput(continueConfig)

		if res == "C" || res == "c" {
			fmt.Printf(
				"\nCongratulations %s! You're cashing out with $%d",
				blackjack.Player.Name,
				blackjack.Player.Cash,
			)
			break
		}

		if blackjack.Player.Cash == 0 {
			fmt.Println("\n*******************************")
			fmt.Println("Busted! You're out of money.")
			break
		}
	}
	fmt.Println("\n*******************************")
	fmt.Println("Thanks for playing! Come back with more cash.")
}

type payoutType int

const (
	normalRate    payoutType = 1
	blackjackRate payoutType = 2
)

// TODO: constructor function for this structure, try to use either Configuration pattern or dependency injection
type Blackjack struct {
	Player  *players.Player
	Dealer  *players.Dealer
	Deck    *decks.BlackjackDeck
	scanner *bufio.Scanner
	// The rate a player's bet is payout at when they win a round
	payoutRate payoutType
}

func (bj *Blackjack) isSplitRound() bool {
	return bj.Player.HasSplitCards()
}

func (bj *Blackjack) DealPlayerCards(count int) {
	cards := bj.Deck.Pop(count)
	for _, c := range cards {
		c.IsFaceUp = true

	}

	bj.Player.Cards = append(bj.Player.Cards, cards...)
}

func (bj *Blackjack) DealDealerCards(count int, isFaceUp bool) {
	cards := bj.Deck.Pop(count)
	for _, c := range cards {
		c.IsFaceUp = isFaceUp
	}

	bj.Dealer.Cards = append(bj.Dealer.Cards, cards...)
}

func (bj *Blackjack) DealFirstCards() {
	fmt.Println("\nDealing cards...")
	// dealt in a loop so each player + dealer is given a card one after the other
	// instead of dealing out a player entirely before moving to the next one
	for i := 0; i < 2; i++ {
		bj.DealPlayerCards(1)

		// dealer card only face up on first card dealt
		if i == 0 {
			bj.DealDealerCards(1, true)
		} else {
			bj.DealDealerCards(1, false)
		}
	}

	bj.PrintTableCards()
}

func (bj *Blackjack) PlaceBet(player *players.Player) {
	fmt.Printf("%s has $%d in wallet\n", player.Name, player.Cash)

	var bet int
	var err error
	inputConfig := utils.NewInputConfig(bj.scanner)

	for {
		fmt.Print("Place your bet: $")
		bet, err = utils.GetUserInputInteger(inputConfig)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("That bet is invalid. Please try again.")
			continue
		}

		if bet > player.Cash {
			fmt.Printf(
				"That bet is larger than your wallet ($%d). Please try again.\n",
				player.Cash,
			)
			continue
		}
		// bet valid
		break
	}

	player.Cash = player.Cash - bet
	player.Bet = bet

	fmt.Printf("Remaining in wallet: $%d\n", player.Cash)
}

func (bj *Blackjack) PrintTableCards() {
	utils.ClearTerminal()
	config := utils.NewPrintTableConfig(bj.Dealer, bj.Player, bj.Deck)
	utils.PrintTable(config)
}

func (bj *Blackjack) PrintSplitRoundCards(round int) {

	utils.ClearTerminal()
	config := utils.NewPrintTableConfig(bj.Dealer, bj.Player, bj.Deck).
		SetTitle("Table Split Round").
		SetSubtitle(fmt.Sprintf("Hand: %d", round))

	utils.PrintTable(config)
}

const (
	// TODO: implement hint machanic that assess the cards on the table and spits out the
	// recommended best move based on basic strategy
	HINT   = "i"
	HIT    = "h"
	STAND  = "s"
	DOUBLE = "d"
	// TODO: implement
	SPLIT = "l"
	// TODO: implement; early and late surrenders? Make it a config options at the beginning screen
	SURRENDER = 'r'
)

// Check if the player has other move options aside from hit and stand
func (bj *Blackjack) GetOtherMoves() string {

	move := ""

	// if original two cards dealt
	if len(bj.Player.Cards) == 2 {

		playerTotal := utils.CalcCardsTotal(bj.Player.Cards)
		if playerTotal == 9 || playerTotal == 10 || playerTotal == 11 {
			move += DOUBLE
		}

		if bj.Player.Cards[0].Rank.Name == bj.Player.Cards[1].Rank.Name && !bj.isSplitRound() {
			// if the first two cards initially dealt are of same name
			// can only split when player has no split cards
			move += SPLIT
		}
	}

	return move
}

func (bj *Blackjack) ChooseNextMove() string {
	var prompt string
	// determine players next available moves
	otherMoves := bj.GetOtherMoves()

	switch otherMoves {
	case DOUBLE:
		prompt = `
	------------------------
	| HIT | STAND | DOUBLE |
	| (h) |  (s)  |  (d)   |
	------------------------
		`
	case SPLIT:
		prompt = `
	-----------------------
	| HIT | STAND | SPLIT |
	| (h) |  (s)  |  (l)  |
	-----------------------
		`

	case DOUBLE + SPLIT:
		prompt = `
	--------------------------------
	| HIT | STAND | DOUBLE | SPLIT |
	| (h) |  (s)  |  (d)   |  (l)  |
	--------------------------------
		`

	default:
		prompt = `
	---------------
	| HIT | STAND |  
	| (h) |  (s)  |
	---------------
		`
	}
	fmt.Println(prompt)

	expectedValues := fmt.Sprintf("hs%s", otherMoves)

	inputConfig := utils.NewInputConfig(bj.scanner).
		SetExpectedValues(strings.Split(expectedValues, "")...)
	var move string
	var err error
	for {
		fmt.Print(PROMPT)
		move, err = utils.GetUserInput(inputConfig)
		if err == nil {
			// input was correct, break out of loop
			break
		}

		fmt.Println(err.Error())
	}
	return move
}

type RoundOutcome int

const (
	InProgress RoundOutcome = 0
	PlayerWon  RoundOutcome = 1
	PlayerLost RoundOutcome = 2
	Standoff   RoundOutcome = 3
	Done       RoundOutcome = 4
)

// Dealer goes through their turn of hitting/standing
// Returns the score of the dealers hand
func (bj *Blackjack) dealersTurn() (dealerScore int) {
	// dealer turns over face down card
	for _, card := range bj.Dealer.Cards {
		card.IsFaceUp = true
	}

	dealerScore = utils.CalcCardsTotal(bj.Dealer.Cards)

	// Check if the dealer has a soft 17 on the first two cards
	// this is an ace + 6 which can either be 7 or 17, so force the
	// dealer to hit assuming the value of the ace will be 1
	if dealerScore == 17 && len(bj.Dealer.Cards) == 2 {
		card := bj.Deck.Pop(1)[0]
		card.IsFaceUp = true
		// deal the dealer one card
		bj.Dealer.Cards = append(bj.Dealer.Cards, card)
		dealerScore = utils.CalcCardsTotal(bj.Dealer.Cards)
	}

	// when the dealer hits a score of 17 or more, auto stand
	for dealerScore < 17 {
		card := bj.Deck.Pop(1)[0]
		card.IsFaceUp = true

		// deal the dealer one card
		bj.Dealer.Cards = append(bj.Dealer.Cards, card)
		// get the new score after the added card
		dealerScore = utils.CalcCardsTotal(bj.Dealer.Cards)
	}
	return dealerScore
}

func (bj *Blackjack) PlayerStand() RoundOutcome {
	dealerScore := bj.dealersTurn()
	bj.PrintTableCards()

	// check if dealer bust
	if dealerScore > BLACKJACK {
		return PlayerWon
	}

	playerScore := utils.CalcCardsTotal(bj.Player.Cards)

	var outcome RoundOutcome
	// compare scores to see who won
	if dealerScore > playerScore {
		outcome = PlayerLost
	} else if dealerScore < playerScore {
		if playerScore == 21 {
			bj.payoutRate = blackjackRate
			fmt.Println("BLACKJACK!!\nCollect your winnings at a rate of 1.5.")
		}
		outcome = PlayerWon
	} else {
		outcome = Standoff
	}

	return outcome
}

// Player was dealt a natural blackjack, check if the dealer also has one
func (bj *Blackjack) DealerBlackjackCheck() RoundOutcome {
	// dealer turns over face down card
	for _, card := range bj.Dealer.Cards {
		card.IsFaceUp = true
	}

	dealerScore := utils.CalcCardsTotal(bj.Dealer.Cards)
	playerScore := utils.CalcCardsTotal(bj.Player.Cards)

	bj.PrintTableCards()

	var outcome RoundOutcome
	// compare scores to see who won
	if dealerScore < playerScore {
		outcome = PlayerWon
	} else {
		outcome = Standoff
	}
	return outcome
}

// Player hit
func (bj *Blackjack) PlayerHit() (outcome RoundOutcome) {
	bj.DealPlayerCards(1)
	bj.PrintTableCards()

	newTotal := utils.CalcCardsTotal(bj.Player.Cards)

	outcome = InProgress

	if newTotal > BLACKJACK {
		// dont stand cause dealer doesnt need to hit
		outcome = PlayerLost
	} else if newTotal == BLACKJACK {
		fmt.Println("BLACKJACK!!\nCollect your winnings at a rate of 1.5.")
		bj.payoutRate = blackjackRate
		outcome = bj.PlayerStand()
	}

	return outcome
}

func (bj *Blackjack) PlayerDouble() (outcome RoundOutcome) {
	// check if the user has enough money to do a double
	if bj.Player.Cash < bj.Player.Bet {
		fmt.Printf(
			"You do not have enough cash left to perform a double.\nNeed $%d, but you only have %d\n",
			bj.Player.Cash,
			bj.Player.Bet,
		)
		return InProgress
	}

	// player has enough cash
	bj.Player.Cash -= bj.Player.Bet
	bj.Player.Bet *= 2
	// only pop one more card before standing
	card := bj.Deck.Pop(1)[0]
	card.IsFaceUp = true

	bj.Player.Cards = append(bj.Player.Cards, card)

	return bj.PlayerStand()
}

func (bj *Blackjack) PlayerSplit() (outlcome RoundOutcome) {
	// check if the user has enough money to do a split
	// need enough cash to double to bet
	if bj.Player.Cash < bj.Player.Bet {
		fmt.Printf(
			"You dont have enough cash to split.\nNeed $%d, but you only have $%d\n",
			bj.Player.Bet*2,
			bj.Player.Bet+bj.Player.Cash,
		)
		return InProgress
	}

	savedBet := bj.Player.Bet

	// cards to be evaluated after both split cards have been finalized by player
	// slice of anonymous structs
	var savedForEvaluation = []struct {
		id  int
		cds []*card.Card
	}{}

	bj.Player.MoveCardsToSplit()

	for idx := 1; bj.Player.HasSplitCards(); idx++ {
		// potential hand to save
		saved := struct {
			id  int
			cds []*card.Card
		}{
			id:  idx,
			cds: []*card.Card{},
		}
		// if the bet was previously cleared pull money from player for next split card bet
		if bj.Player.Bet == 0 {
			bj.Player.Bet = savedBet
			bj.Player.Cash -= savedBet
		}

		// move split card to players cards
		bj.Player.NextSplitCard()

		if bj.Player.Cards[0].Rank.Name == rank.Ace {
			// ace cards only allowed one card on split
			// deal card and add to pending array

			bj.DealPlayerCards(1)
			bj.PrintSplitRoundCards(saved.id)

			// prompt user to hit any key to continue
			utils.EnterToContinue(bj.scanner)

			saved.cds = append(saved.cds, bj.Player.Cards...)
			savedForEvaluation = append(savedForEvaluation, saved)
			// assign new slice instead of reslice because need that slice in memory to maintain
			// the same values so savedForEvaluation can be looped on later with the same values
			// otherwise if reslice to {:0] it would overwrite the values
			bj.Player.Cards = []*card.Card{}
			continue
		}

		// draw one card and add to player cards
		bj.DealPlayerCards(1)
		bj.PrintSplitRoundCards(saved.id)

		localOutcome := InProgress
		for localOutcome == InProgress {
			move := bj.ChooseNextMove()
			switch move {
			case HIT:
				// custom hit logic for split rounds
				bj.DealPlayerCards(1)
				bj.PrintSplitRoundCards(saved.id)

				newTotal := utils.CalcCardsTotal(bj.Player.Cards)

				if newTotal > BLACKJACK {
					// dont stand cause dealer doesnt need to hit
					bj.PlayerLostHand()
					bj.Deck.AddDiscardedCards(bj.Player.Cards)
					bj.Player.Cards = []*card.Card{}
					localOutcome = Done
					utils.EnterToContinue(bj.scanner)
					continue
				} else if newTotal == BLACKJACK {

					// theres no blackjack rate during a split
					saved.cds = append(saved.cds, bj.Player.Cards...)
					savedForEvaluation = append(savedForEvaluation, saved)
					bj.Player.Cards = []*card.Card{}
					fmt.Println("Player hit max card score.")
					localOutcome = Done
					utils.EnterToContinue(bj.scanner)
				}
			case STAND:
				saved.cds = append(saved.cds, bj.Player.Cards...)
				// custom stand logic for split rounds
				savedForEvaluation = append(savedForEvaluation, saved)

				bj.Player.Cards = []*card.Card{}
				localOutcome = Done
			case DOUBLE:
				// TODO: implement split double
				// outcome = bj.PlayerDouble()
				fmt.Println("Split double under construction. Try something else.")
				localOutcome = InProgress
			}
		}
	}

	dealerTotal := bj.dealersTurn()
	for idx, savedStruct := range savedForEvaluation {

		if bj.Player.Bet == 0 {
			bj.Player.Bet = savedBet
			bj.Player.Cash -= savedBet
		}

		bj.Player.Cards = savedStruct.cds
		bj.PrintSplitRoundCards(savedStruct.id)

		playerTotal := utils.CalcCardsTotal(savedStruct.cds)

		if playerTotal > dealerTotal || dealerTotal > BLACKJACK {
			bj.PlayerWonHand()
		} else if dealerTotal > playerTotal {
			bj.PlayerLostHand()
		} else if dealerTotal == playerTotal {
			bj.Standoff()
		}

		if idx != len(savedForEvaluation)-1 {
			// don't show on the last item cause they'll see it anyways once split is finalized
			// pause for user to look at at outcome
			utils.EnterToContinue(bj.scanner)
		}
	}

	bj.cleanup()
	return Done
}

// Player has lost the hand, clean up for next deal
func (bj *Blackjack) PlayerLostHand() {
	fmt.Print("***  Dealer win!  ***\nCollecting all losing bets...\n\n")
	// Player loses the bet
	bj.Player.Bet = 0
}

func (bj *Blackjack) PlayerWonHand() {
	fmt.Print("***  Player win!  ***\nAdding winnings to your wallet...\n\n")
	if bj.payoutRate == normalRate {
		// effective rate or 1
		bj.Player.Cash += (bj.Player.Bet * 2)
	} else {
		// else is blackjack rate so rate is 1.5
		// *15 /10 is the same as doing * 1.5 but we avoid floats
		bj.Player.Cash += (bj.Player.Bet * 15 / 10)
	}
	bj.Player.Bet = 0
}

// Player and dealer have the same card total
func (bj *Blackjack) Standoff() {
	fmt.Print("Push! Returning all bets...\n\n")
	// player looses nothing, add bet back to cash
	bj.Player.Cash += bj.Player.Bet
	bj.Player.Bet = 0
}

// Clear the cards of the dealer and the player by clearing and reslicing
// which will maintain the same memory adderss for the game, but also
// retain the highest capacity acheived across hands played
func (bj *Blackjack) cleanup() {
	bj.Deck.AddDiscardedCards(bj.Player.Cards)
	// sets values in slice to 'zero'
	clear(bj.Player.Cards)
	// reslice so length of slice is now zero again
	bj.Player.Cards = bj.Player.Cards[:0]

	// same as above
	bj.Deck.AddDiscardedCards(bj.Dealer.Cards)
	clear(bj.Dealer.Cards)
	bj.Dealer.Cards = bj.Dealer.Cards[:0]

	bj.payoutRate = normalRate
}
