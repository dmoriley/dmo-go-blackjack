package game

import (
	"blackjack/card"
	"blackjack/card/rank"
	"blackjack/deck"
	"blackjack/game/players"
	"blackjack/game/utils"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	// Max value to win or bust
	BLACKJACK = 21
	PROMPT    = "=> "
)

type Blackjack struct {
	Player *players.Player
	Dealer *players.Dealer
	Deck   *BlackjackDeck
}

func getUserInput(scanner *bufio.Scanner) (string, error) {

	scanned := scanner.Scan()
	if !scanned {
		return "", scanner.Err()
	}

	return strings.TrimSpace(scanner.Text()), nil
}

// Parse user intput for a number
func getUserInputInteger(scanner *bufio.Scanner) (int, error) {
	text, err := getUserInput(scanner)
	if err != nil {
		return 0, err
	}

	integer, err := strconv.Atoi(text)
	if err != nil {
		return 0, err
	}
	return integer, nil
}

// Start the game
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	player := &players.Player{
		Name: "David",
		Cash: 500,
	}
	dealer := &players.Dealer{}
	deck := NewBlackjackDeck(3)
	deck.Shuffle(5)

	blackjack := &Blackjack{
		Dealer: dealer,
		Player: player,
		Deck:   deck,
	}

	PlaceBet(blackjack.Player, scanner)
	DealInitialCards(blackjack)

	for {
		fmt.Fprintf(out, PROMPT)
		input, _ := getUserInput(scanner)
		io.WriteString(out, input+"\n")
	}
}

// Determine the card value total of the hand supplied
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

func DealInitialCards(bj *Blackjack) {
	fmt.Println("\nDealing initial cards...")
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

	fmt.Println(PrintTableCards(bj))
}

func PlaceBet(player *players.Player, scanner *bufio.Scanner) {
	prompt := fmt.Sprintf("%s has $%d in wallet", player.Name, player.Cash)
	fmt.Println(prompt)

	var bet int
	var err error

	for {
		fmt.Print("Place your bet: $")
		bet, err = getUserInputInteger(scanner)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("That bet is invalid. Please try again.")
			continue
		}

		if bet > player.Cash {
			prompt = fmt.Sprintf(
				"That bet is larger than your wallet ($%d). Please try again.",
				player.Cash,
			)
			fmt.Println(prompt)
			continue
		}
		// bet valid
		break
	}

	player.Cash = player.Cash - bet
	player.Bet = bet

	prompt = fmt.Sprintf("You have place a bet of: $%d\nRemaining in wallet: $%d", bet, player.Cash)
	fmt.Println(prompt)
}

func PrintTableCards(bj *Blackjack) string {
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
	out.WriteString(deck.PrettyPrintCards(bj.Dealer.Cards))
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
	out.WriteString(deck.PrettyPrintCards(bj.Player.Cards))
	utils.FillTextAndPad(&out, 45, ' ', '*', "", "")
	utils.FillTextAndPad(&out, 45, '*', '*', "", "")
	utils.FillTextAndPad(&out, 45, '*', '*', "", "")

	return out.String()
}
