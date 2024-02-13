package utils

import (
	"blackjack/decks"
	"blackjack/game/players"
	"bytes"
	"fmt"
)

const (
	TABLE_CHAR_WIDTH = 45
)

// configuration pattern
type printTableConfig struct {
	dealer   *players.Dealer
	player   *players.Player
	deck     *decks.BlackjackDeck
	title    string
	subtitle string
}

func NewPrintTableConfig(
	dealer *players.Dealer,
	player *players.Player,
	deck *decks.BlackjackDeck,
) *printTableConfig {

	return &printTableConfig{
		dealer:   dealer,
		player:   player,
		deck:     deck,
		title:    "Table Cards",
		subtitle: "",
	}
}

func (c *printTableConfig) SetTitle(title string) *printTableConfig {
	c.title = title
	return c
}

func (c *printTableConfig) SetSubtitle(sub string) *printTableConfig {
	c.subtitle = sub
	return c
}

func PrintTable(config *printTableConfig) {
	var out bytes.Buffer

	out.WriteString("\n")
	FillTextAndPad(&out, TABLE_CHAR_WIDTH, '*', '*', "", "")
	FillTextAndPad(&out, TABLE_CHAR_WIDTH, '*', '*', config.title, "middle")
	if len(config.subtitle) > 0 {
		FillTextAndPad(&out, TABLE_CHAR_WIDTH, '*', '*', config.subtitle, "right")
	}
	FillTextAndPad(&out, TABLE_CHAR_WIDTH, '*', '*', "", "")
	FillTextAndPad(&out, TABLE_CHAR_WIDTH, ' ', '*', "", "")
	FillTextAndPad(&out, TABLE_CHAR_WIDTH, ' ', '*', "Dealers cards", "left")
	FillTextAndPad(
		&out,
		TABLE_CHAR_WIDTH,
		' ',
		'*',
		fmt.Sprintf("Total: %d", CalcCardsTotal(config.dealer.Cards)),
		"left",
	)
	FillTextAndPad(&out, TABLE_CHAR_WIDTH, ' ', '*', "-------------", "left")
	FillTextAndPad(&out, TABLE_CHAR_WIDTH, ' ', '*', "", "")
	out.WriteString(decks.PrettyPrintCards(config.dealer.Cards))
	FillTextAndPad(&out, TABLE_CHAR_WIDTH, ' ', '*', "", "")

	// player name and card total
	FillTextAndPad(
		&out,
		TABLE_CHAR_WIDTH,
		' ',
		'*',
		fmt.Sprintf("%s cards", config.player.Name),
		"left",
	)
	FillTextAndPad(
		&out,
		TABLE_CHAR_WIDTH,
		' ',
		'*',
		fmt.Sprintf("Total: %d", CalcCardsTotal(config.player.Cards)),
		"left",
	)

	FillTextAndPad(&out, TABLE_CHAR_WIDTH, ' ', '*', "-------------", "left")
	FillTextAndPad(&out, TABLE_CHAR_WIDTH, ' ', '*', "", "")
	out.WriteString(decks.PrettyPrintCards(config.player.Cards))
	FillTextAndPad(&out, TABLE_CHAR_WIDTH, ' ', '*', "", "")
	FillTextAndPad(&out, TABLE_CHAR_WIDTH, '*', '*', "", "")
	FillTextAndPad(
		&out,
		TABLE_CHAR_WIDTH,
		'*',
		'*',
		fmt.Sprintf("%d/%d", config.deck.GetLength(), config.deck.DeckCount*52),
		"middle",
	)
	FillTextAndPad(&out, TABLE_CHAR_WIDTH, '*', '*', "", "")

	fmt.Println(out.String())
}
