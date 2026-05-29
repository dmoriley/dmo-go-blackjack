package utils

import (
	"bytes"
	"fmt"
)

const (
	TABLE_CHAR_WIDTH = 45
)

// configuration pattern

type TableCard struct {
	Label  string
	Value  int
	FaceUp bool
}

func (c TableCard) Inspect() string {
	if !c.FaceUp {
		return "{Face down}"
	}

	return fmt.Sprintf("{%s, value: %d}", c.Label, c.Value)
}

type TableHand struct {
	Label string
	Total int
	Cards []TableCard
}

type printTableConfig struct {
	dealer        TableHand
	player        TableHand
	deckRemaining int
	deckTotal     int
	title         string
	subtitle      string
}

func NewPrintTableConfig(
	dealer TableHand,
	player TableHand,
	deckRemaining int,
	deckTotal int,
) *printTableConfig {

	return &printTableConfig{
		dealer:        dealer,
		player:        player,
		deckRemaining: deckRemaining,
		deckTotal:     deckTotal,
		title:         "Table Cards",
		subtitle:      "",
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
	fmt.Print(RenderTable(config))
}

func RenderTable(config *printTableConfig) string {
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
		fmt.Sprintf("Total: %d", config.dealer.Total),
		"left",
	)
	FillTextAndPad(&out, TABLE_CHAR_WIDTH, ' ', '*', "-------------", "left")
	FillTextAndPad(&out, TABLE_CHAR_WIDTH, ' ', '*', "", "")
	out.WriteString(prettyPrintCards(config.dealer.Cards))
	FillTextAndPad(&out, TABLE_CHAR_WIDTH, ' ', '*', "", "")

	// player name and card total
	FillTextAndPad(
		&out,
		TABLE_CHAR_WIDTH,
		' ',
		'*',
		fmt.Sprintf("%s cards", config.player.Label),
		"left",
	)
	FillTextAndPad(
		&out,
		TABLE_CHAR_WIDTH,
		' ',
		'*',
		fmt.Sprintf("Total: %d", config.player.Total),
		"left",
	)

	FillTextAndPad(&out, TABLE_CHAR_WIDTH, ' ', '*', "-------------", "left")
	FillTextAndPad(&out, TABLE_CHAR_WIDTH, ' ', '*', "", "")
	out.WriteString(prettyPrintCards(config.player.Cards))
	FillTextAndPad(&out, TABLE_CHAR_WIDTH, ' ', '*', "", "")
	FillTextAndPad(&out, TABLE_CHAR_WIDTH, '*', '*', "", "")
	FillTextAndPad(
		&out,
		TABLE_CHAR_WIDTH,
		'*',
		'*',
		fmt.Sprintf("%d/%d", config.deckRemaining, config.deckTotal),
		"middle",
	)
	FillTextAndPad(&out, TABLE_CHAR_WIDTH, '*', '*', "", "")

	return out.String()
}

func prettyPrintCards(cards []TableCard) string {
	var out bytes.Buffer

	if len(cards) == 0 {
		out.WriteString("{ No Cards }")
		return out.String()
	}

	out.WriteString("{\n")
	for _, card := range cards {
		out.WriteString(fmt.Sprintf("\t%s\n", card.Inspect()))
	}
	out.WriteString("}\n")
	return out.String()
}
