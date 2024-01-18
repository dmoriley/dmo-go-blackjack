package game

import (
	"blackjack/card"
	"blackjack/card/rank"
	"blackjack/decks"
	"blackjack/game/players"
	"blackjack/game/utils"
	"bufio"
	"bytes"
	"fmt"
	"io"
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
	config := decks.NewBlackjackDeckConfig().WithNumberOfDecks(6)
	deck := decks.NewBlackjackDeck(config)
	deck.Shuffle(5)

	blackjack := &Blackjack{
		Dealer:  dealer,
		Player:  player,
		Deck:    deck,
		scanner: scanner,
	}

	for {
		blackjack.PlaceBet(blackjack.Player)
		blackjack.DealCards()
		var roundOver bool
		if playerScore := GetCardsTotal(blackjack.Player.Cards); playerScore >= 21 {
			roundOver = true
			// player hit 21 points or is over, in either case
			// they have no more moves to play, so stand
			if playerScore == 21 {
				blackjack.DealerBlackjackCheck()
			} else {
				blackjack.PlayerStand()
			}
		}
		for !roundOver {
			move := blackjack.ChooseNextMove()

			switch move {
			case "h":
				roundOver = blackjack.PlayerHit()
			case "s":
				blackjack.PlayerStand()
				roundOver = true
			}
		}

		if blackjack.Player.Cash == 0 {
			fmt.Println("\n*******************************")
			fmt.Println("Busted! You're out of money.")
			break
		}
	}
	fmt.Println("*******************************")
	fmt.Println("Thanks for playing! Come back with more cash.")
}

// TODO: constructor function for this structure, try to use either Configuration pattern or dependency injection
type Blackjack struct {
	Player  *players.Player
	Dealer  *players.Dealer
	Deck    *decks.BlackjackDeck
	scanner *bufio.Scanner
}

func (bj *Blackjack) DealCards() {
	fmt.Println("\nDealing cards...")
	// dealt in a loop so each player + dealer is given a card one after the other
	// instead of dealing out a player entirely before moving to the next one
	for i := 0; i < 2; i++ {
		card := bj.Deck.Pop(1)[0]
		card.IsFaceUp = true

		bj.Player.Cards = append(bj.Player.Cards, card)

		card = bj.Deck.Pop(1)[0]
		// dealer card only face up on first card dealt
		if i == 0 {
			card.IsFaceUp = true
		}
		bj.Dealer.Cards = append(bj.Dealer.Cards, card)
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
	var out bytes.Buffer

	out.WriteString("\n")
	utils.FillTextAndPad(&out, 45, '*', '*', "", "")
	utils.FillTextAndPad(&out, 45, '*', '*', "Table Cards", "middle")
	utils.FillTextAndPad(&out, 45, '*', '*', "", "")
	utils.FillTextAndPad(&out, 45, ' ', '*', "", "")
	utils.FillTextAndPad(&out, 45, ' ', '*', "Dealers cards", "left")
	utils.FillTextAndPad(
		&out,
		45,
		' ',
		'*',
		fmt.Sprintf("Total: %d", GetCardsTotal(bj.Dealer.Cards)),
		"left",
	)
	utils.FillTextAndPad(&out, 45, ' ', '*', "-------------", "left")
	utils.FillTextAndPad(&out, 45, ' ', '*', "", "")
	out.WriteString(decks.PrettyPrintCards(bj.Dealer.Cards))
	utils.FillTextAndPad(&out, 45, ' ', '*', "", "")

	// player name and card total
	utils.FillTextAndPad(&out, 45, ' ', '*', fmt.Sprintf("%s cards", bj.Player.Name), "left")
	utils.FillTextAndPad(
		&out,
		45,
		' ',
		'*',
		fmt.Sprintf("Total: %d", GetCardsTotal(bj.Player.Cards)),
		"left",
	)

	utils.FillTextAndPad(&out, 45, ' ', '*', "-------------", "left")
	utils.FillTextAndPad(&out, 45, ' ', '*', "", "")
	out.WriteString(decks.PrettyPrintCards(bj.Player.Cards))
	utils.FillTextAndPad(&out, 45, ' ', '*', "", "")
	utils.FillTextAndPad(&out, 45, '*', '*', "", "")
	utils.FillTextAndPad(&out, 45, '*', '*', "", "")

	fmt.Println(out.String())
}

// Determine the card value total of the hand supplied
// Does not count value of card that is not face up
func GetCardsTotal(cards []*card.Card) int {
	total := 0
	aceCount := 0

	for _, card := range cards {
		if !card.IsFaceUp {
			// dont total cards that arent being shown
			continue
		}
		switch card.Rank.Name {
		case rank.Ace:
			aceCount++
		default:
			total += card.Rank.Value
		}
	}

	if aceCount == 1 {
		if total+11 <= BLACKJACK {
			total += 11
		} else {
			total += 1
		}
	} else if aceCount > 1 {
		// check if one ace can have the value of 11
		if total+11+(aceCount-1) <= BLACKJACK {
			total += (11 + (aceCount - 1))
		} else {
			// all aces must have a value of one
			total += aceCount
		}
	}

	return total
}

type PlayerMove int

const (
	// TODO: implement hint machanic that assess the cards on the table and spits out the
	// recommended best move based on basic strategy
	HINT  PlayerMove = 0
	HIT   PlayerMove = 1
	STAND PlayerMove = 2
	// TODO: implement
	DOUBLE PlayerMove = 3
	// TODO: implement
	SPLIT PlayerMove = 4
	// TODO: implement
	SURRENDER PlayerMove = 5
)

func (bj *Blackjack) GetAvailableMoves() PlayerMove {
	// hs for hit and stand
	return HIT + STAND
}

func (bj *Blackjack) ChooseNextMove() string {
	var prompt string
	// determine players next available moves
	availableMoves := bj.GetAvailableMoves()
	switch availableMoves {
	case HIT + STAND:
		prompt = `
	---------------
	| HIT | STAND |  
	| (h) |  (s)  |
	---------------
		`
	}
	fmt.Println(prompt)

	inputConfig := utils.NewInputConfig(bj.scanner).SetExpectedValues("h", "s")
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

// TODO: refactor method to return the outcome of the stand instead
// of calling methods inside
func (bj *Blackjack) PlayerStand() {
	// dealer turns over face down card
	for _, card := range bj.Dealer.Cards {
		card.IsFaceUp = true
	}

	dealerScore := GetCardsTotal(bj.Dealer.Cards)

	// Check if the dealer has a soft 17 on the first two cards
	// this is an ace + 6 which can either be 7 or 17, so force the
	// dealer to hit assuming the value of the ace will be 1
	if dealerScore == 17 && len(bj.Dealer.Cards) == 2 {
		card := bj.Deck.Pop(1)[0]
		card.IsFaceUp = true
		// deal the dealer one card
		bj.Dealer.Cards = append(bj.Dealer.Cards, card)
		dealerScore = GetCardsTotal(bj.Dealer.Cards)
	}

	// when the dealer hits a score of 17 or more, auto stand
	for dealerScore < 17 {
		card := bj.Deck.Pop(1)[0]
		card.IsFaceUp = true

		// deal the dealer one card
		bj.Dealer.Cards = append(bj.Dealer.Cards, card)
		// get the new score after the added card
		dealerScore = GetCardsTotal(bj.Dealer.Cards)
	}

	bj.PrintTableCards()

	// check if dealer bust
	if dealerScore > BLACKJACK {
		bj.PlayerWonHand()
		return
	}

	playerScore := GetCardsTotal(bj.Player.Cards)

	// compare scores to see who won
	if dealerScore > playerScore {
		bj.PlayerLostHand()
	} else if dealerScore < playerScore {
		bj.PlayerWonHand()
	} else {
		bj.Standoff()
	}

}

// Player was dealt a natural blackjack, check if the dealer also has one
func (bj *Blackjack) DealerBlackjackCheck() {
	// dealer turns over face down card
	for _, card := range bj.Dealer.Cards {
		card.IsFaceUp = true
	}

	dealerScore := GetCardsTotal(bj.Dealer.Cards)
	playerScore := GetCardsTotal(bj.Player.Cards)

	bj.PrintTableCards()

	// compare scores to see who won
	if dealerScore < playerScore {
		bj.PlayerWonHand()
	} else {
		bj.Standoff()
	}

}

// Player hit
func (bj *Blackjack) PlayerHit() (roundOver bool) {
	card := bj.Deck.Pop(1)[0]
	card.IsFaceUp = true

	bj.Player.Cards = append(bj.Player.Cards, card)

	bj.PrintTableCards()

	newTotal := GetCardsTotal(bj.Player.Cards)

	if newTotal > BLACKJACK {
		// dont stand cause dealer doesnt need to hit
		bj.PlayerLostHand()
		roundOver = true
	} else if newTotal == BLACKJACK {
		bj.PlayerStand()
		roundOver = true
	}

	return
}

// Player has lost the hand, clean up for next deal
func (bj *Blackjack) PlayerLostHand() {
	fmt.Print("***  Dealer win!  ***\nCollecting all losing bets...\n\n")
	// Player loses the bet
	bj.Player.Bet = 0
	bj.cleanupCards()
}

func (bj *Blackjack) PlayerWonHand() {
	fmt.Print("***  Player win!  ***\nAdding winnings to your wallet...\n\n")
	winnings := bj.Player.Bet * 2
	bj.Player.Cash += winnings
	bj.Player.Bet = 0
	bj.cleanupCards()
}

// Player and dealer have the same card total
func (bj *Blackjack) Standoff() {
	fmt.Print("Push! Returning all bets...\n\n")
	// player looses nothing, add bet back to cash
	bj.Player.Cash += bj.Player.Bet
	bj.Player.Bet = 0
	bj.cleanupCards()
}

// Clear the cards of the dealer and the player by clearing and reslicing
// which will maintain the same memory adderss for the game, but also
// retain the highest capacity acheived across hands played
func (bj *Blackjack) cleanupCards() {
	bj.Deck.AddDiscardedCards(bj.Player.Cards)
	// sets values in slice to 'zero'
	clear(bj.Player.Cards)
	// reslice so length of slice is now zero again
	bj.Player.Cards = bj.Player.Cards[:0]
	// same as above
	bj.Deck.AddDiscardedCards(bj.Dealer.Cards)
	clear(bj.Dealer.Cards)
	bj.Dealer.Cards = bj.Dealer.Cards[:0]
}
