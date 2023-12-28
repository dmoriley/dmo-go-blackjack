package game

import (
	"blackjack/card"
	"blackjack/card/rank"
	"blackjack/deck"
	"blackjack/game/players"
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	// Max value to win or bust
	BLACKJACK = 21
	PROMPT    = "|= "
)

type Blackjack struct {
	Player *players.Player
	Dealer *players.Dealer
	Deck   *deck.Deck
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
	david := &players.Player{
		Name: "David",
		Cash: 500,
	}

	for {
		fmt.Fprintf(out, PROMPT)
		input, _ := getUserInput(scanner)
		io.WriteString(out, input+"\n")
		PlaceBet(david, scanner)
	}
}

// Determine the card value total of the hand supplied
func GetCardsTotal(cards []*card.Card) int {
	total := 0
	aceCount := 0

	for _, card := range cards {
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
