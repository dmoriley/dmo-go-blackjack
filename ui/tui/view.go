package tui

import (
	"fmt"
	"strings"

	"blackjack/game"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	content := ""
	switch m.screen {
	case screenNameEntry:
		content = m.viewNameEntry()
	case screenEnd:
		content = m.viewEnd()
	default:
		content = m.viewTable()
	}

	return m.styles.app.Render(content)
}

func (m Model) viewNameEntry() string {
	sections := []string{
		m.styles.headerTitle.Render("Blackjack"),
		m.styles.muted.Render("A restrained table view built on the headless engine."),
		"",
		m.styles.namePrompt.Render("Choose the name that will sit at the table."),
		m.styles.input.Render(m.nameInput.View()),
		m.renderError(),
		m.styles.muted.Render("Press enter to start. Press esc to quit."),
	}

	card := m.styles.nameShell.Render(strings.Join(filterEmpty(sections), "\n"))
	placed := lipgloss.Place(
		m.width-4,
		m.height-4,
		lipgloss.Center,
		lipgloss.Center,
		card,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(m.styles.canvas.GetForeground()),
		lipgloss.WithWhitespaceBackground(m.styles.canvas.GetBackground()),
	)
	return m.styles.canvas.Render(placed)
}

func (m Model) viewTable() string {
	header := m.viewHeader()
	table := m.viewMainTable()
	sidebar := m.viewSidebar()
	body := lipgloss.JoinHorizontal(lipgloss.Top, table, "  ", sidebar)
	footer := m.viewFooter()

	sections := []string{header, body, footer}
	return m.styles.shell.Render(strings.Join(filterEmpty(sections), "\n\n"))
}

func (m Model) viewHeader() string {
	phase := phaseLabel(m.snapshot.Phase)
	pill := m.styles.statusPill.Render(phase)
	if m.snapshot.Phase == game.PhaseRoundResult || m.snapshot.Phase == game.PhaseGameOver || m.snapshot.Phase == game.PhaseCashedOut {
		pill = m.styles.statusPillAlert.Render(phase)
	}

	left := m.styles.headerTitle.Render("Blackjack")
	right := lipgloss.JoinHorizontal(lipgloss.Left,
		pill,
		"  ",
		m.styles.headerMeta.Render(fmt.Sprintf("Player %s", m.playerName)),
		"  ",
		m.styles.headerMeta.Render(fmt.Sprintf("Bankroll $%d", m.snapshot.Cash)),
	)

	line := lipgloss.JoinHorizontal(lipgloss.Top, left, strings.Repeat(" ", max(2, m.width-lipgloss.Width(left)-lipgloss.Width(right)-10)), right)
	return m.styles.header.Render(line)
}

func (m Model) viewMainTable() string {
	dealer := m.renderHandPanel("Dealer", m.snapshot.Dealer, false, false, 0)
	hands := []string{dealer}

	for idx, hand := range m.snapshot.PlayerHands {
		label := m.snapshot.PlayerName
		if m.snapshot.IsSplitRound {
			label = fmt.Sprintf("%s · Hand %d", m.snapshot.PlayerName, idx+1)
		}
		active := idx == m.snapshot.ActiveHandIndex && m.snapshot.Phase == game.PhasePlayerTurn
		hands = append(hands, m.renderHandPanel(label, hand, true, active, idx))
	}

	return lipgloss.JoinVertical(lipgloss.Left, hands...)
}

func (m Model) renderHandPanel(title string, hand game.HandState, isPlayer bool, active bool, idx int) string {
	style := m.styles.inactiveHand
	if active {
		style = m.styles.activeHand
	}

	meta := []string{fmt.Sprintf("Total %d", hand.Total)}
	if isPlayer && hand.Bet > 0 {
		meta = append(meta, fmt.Sprintf("Bet $%d", hand.Bet))
	}
	if hand.Blackjack {
		meta = append(meta, "Blackjack")
	}
	if hand.Busted {
		meta = append(meta, "Busted")
	}
	if hand.Resolved {
		meta = append(meta, fmt.Sprintf("Result %s", outcomeLabel(hand.Outcome)))
	}

	header := lipgloss.JoinHorizontal(lipgloss.Top,
		m.styles.panelTitle.Render(title),
		"  ",
		m.styles.muted.Render(strings.Join(meta, " · ")),
	)

	return style.Render(lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		m.renderCards(hand.Cards),
	))
}

func (m Model) renderCards(cards []game.CardState) string {
	if len(cards) == 0 {
		return m.styles.muted.Render("No cards in play")
	}

	rendered := make([]string, 0, len(cards))
	for _, card := range cards {
		rendered = append(rendered, m.renderCard(card))
	}

	rows := make([]string, 0, (len(rendered)+3)/4)
	for start := 0; start < len(rendered); start += 4 {
		end := min(start+4, len(rendered))
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, rendered[start:end]...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m Model) renderCard(card game.CardState) string {
	if !card.FaceUp {
		return m.styles.cardFaceDown.Render(strings.Join(faceDownCardLines(), "\n"))
	}

	return m.styles.card.Render(strings.Join(faceUpCardLines(card), "\n"))
}

func (m Model) viewSidebar() string {
	status := m.styles.panel.Render(strings.Join(filterEmpty([]string{
		m.styles.panelTitle.Render("Status"),
		fmt.Sprintf("Cash $%d", m.snapshot.Cash),
		fmt.Sprintf("Previous bet $%d", m.snapshot.PreviousBet),
		fmt.Sprintf("Deck %d / %d", m.snapshot.DeckRemaining, m.snapshot.DeckTotal),
		fmt.Sprintf("Pending results %d", m.snapshot.PendingResults),
		m.activeHandSummary(),
	}), "\n"))

	logLines := make([]string, 0, len(m.logs)+1)
	logLines = append(logLines, m.styles.panelTitle.Render("Table Log"))
	if len(m.logs) == 0 {
		logLines = append(logLines, m.styles.muted.Render("The dealer is waiting for the next action."))
	} else {
		for _, entry := range m.logs {
			logLines = append(logLines, m.renderLogEntry(entry))
		}
	}

	logPanel := m.styles.logPanel.Render(strings.Join(logLines, "\n"))
	return lipgloss.JoinVertical(lipgloss.Left, status, "", logPanel)
}

func (m Model) activeHandSummary() string {
	if len(m.snapshot.PlayerHands) == 0 || m.snapshot.ActiveHandIndex < 0 || m.snapshot.ActiveHandIndex >= len(m.snapshot.PlayerHands) {
		return m.styles.muted.Render("No active hand")
	}

	hand := m.snapshot.PlayerHands[m.snapshot.ActiveHandIndex]
	label := "Active hand"
	if m.snapshot.IsSplitRound {
		label = fmt.Sprintf("Active hand %d", m.snapshot.ActiveHandIndex+1)
	}
	return fmt.Sprintf("%s total %d", label, hand.Total)
}

func (m Model) renderLogEntry(entry logEntry) string {
	style := m.styles.muted
	switch entry.tone {
	case logToneSuccess:
		style = m.styles.success
	case logToneWarning:
		style = m.styles.warning
	case logToneDanger:
		style = m.styles.danger
	case logToneInfo:
		style = m.styles.muted
	}

	return style.Render("• " + entry.text)
}

func (m Model) viewFooter() string {
	actions := m.renderActions()
	help := []string{actions}
	if msg := m.renderError(); msg != "" {
		help = append(help, msg)
	}
	return m.styles.actionBar.Render(strings.Join(help, "\n"))
}

func (m Model) renderActions() string {
	switch m.snapshot.Phase {
	case game.PhaseBetting:
		input := lipgloss.JoinHorizontal(lipgloss.Center,
			m.styles.inputLabel.Render("Bet"),
			" ",
			m.styles.input.Render(m.betInput.View()),
		)
		actions := []string{
			m.renderAction("enter", "place bet", true),
			m.renderAction(m.keys.previousBet, "repeat previous bet", m.snapshot.PreviousBet > 0),
			m.renderAction(m.keys.cashOut, "cash out", true),
		}
		return lipgloss.JoinVertical(lipgloss.Left,
			input,
			"",
			lipgloss.JoinHorizontal(lipgloss.Left, actions...),
		)
	case game.PhasePlayerTurn:
		actions := []string{
			m.renderAction(m.keys.hit, "hit", m.hasLegalMove(game.MoveHit)),
			m.renderAction(m.keys.stand, "stand", m.hasLegalMove(game.MoveStand)),
			m.renderAction(m.keys.double, "double", m.hasLegalMove(game.MoveDouble)),
			m.renderAction(m.keys.split, "split", m.hasLegalMove(game.MoveSplit)),
		}
		return lipgloss.JoinHorizontal(lipgloss.Left, actions...)
	case game.PhaseRoundResult:
		return m.renderAction("enter", "continue", true)
	default:
		return m.styles.muted.Render("The hand is complete.")
	}
}

func (m Model) renderAction(key string, label string, enabled bool) string {
	style := m.styles.actionDisabled
	if enabled {
		style = m.styles.muted
	}
	text := fmt.Sprintf("%s %s", m.styles.actionKey.Render(key), label)
	return style.Render(text + "   ")
}

func (m Model) renderError() string {
	if m.errMsg == "" {
		return ""
	}
	return m.styles.warning.Render(m.errMsg)
}

func (m Model) viewEnd() string {
	title := "Session Complete"
	subtitle := fmt.Sprintf("%s leaves the table with $%d.", m.playerName, m.snapshot.Cash)
	if m.endPhase == game.PhaseGameOver {
		title = "Game Over"
		subtitle = fmt.Sprintf("%s is out of cash.", m.playerName)
	}

	body := []string{
		m.styles.headerTitle.Render(title),
		subtitle,
		"",
		m.styles.muted.Render("Press n to start a new game with the same name."),
		m.styles.muted.Render("Press q to leave the table."),
	}

	if len(m.logs) > 0 {
		body = append(body, "")
		body = append(body, m.styles.panelTitle.Render("Final Notes"))
		for _, entry := range m.logs {
			body = append(body, m.renderLogEntry(entry))
		}
	}

	panel := m.styles.endPanel.Render(strings.Join(body, "\n"))
	placed := lipgloss.Place(
		m.width-4,
		m.height-4,
		lipgloss.Center,
		lipgloss.Center,
		panel,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(m.styles.canvas.GetForeground()),
		lipgloss.WithWhitespaceBackground(m.styles.canvas.GetBackground()),
	)
	return m.styles.canvas.Render(placed)
}

func phaseLabel(phase game.Phase) string {
	switch phase {
	case game.PhaseBetting:
		return "Betting"
	case game.PhasePlayerTurn:
		return "Player Turn"
	case game.PhaseRoundResult:
		return "Round Result"
	case game.PhaseGameOver:
		return "Game Over"
	case game.PhaseCashedOut:
		return "Cashed Out"
	default:
		return string(phase)
	}
}

func outcomeLabel(outcome game.Outcome) string {
	switch outcome {
	case game.OutcomeWon:
		return "Won"
	case game.OutcomeLost:
		return "Lost"
	case game.OutcomePush:
		return "Push"
	default:
		return string(outcome)
	}
}

func filterEmpty(values []string) []string {
	filtered := make([]string, 0, len(values))
	for _, value := range values {
		if value != "" {
			filtered = append(filtered, value)
		}
	}
	return filtered
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func faceUpCardLines(card game.CardState) []string {
	rank := shortRank(card.Rank)
	suit := suitSymbol(card.Suit)
	rows := cardInteriorRows(rank, suit)
	lines := make([]string, 0, len(rows)+3)
	lines = append(lines, "┌───────────┐")
	for _, row := range rows {
		lines = append(lines, "│"+row+"│")
	}
	lines = append(lines, "└───────────┘")
	lines = append(lines, fmt.Sprintf("%s of %s [%d]", card.Rank, card.Suit, card.Value))
	return lines
}

func faceDownCardLines() []string {
	return []string{
		"┌───────────┐",
		"│░░░░░░░░░░░│",
		"│░ ▒▒▒▒▒▒▒ ░│",
		"│░ ▒ ┌─┐ ▒ ░│",
		"│░ ▒ │░│ ▒ ░│",
		"│░ ▒ └─┘ ▒ ░│",
		"│░ ▒▒▒▒▒▒▒ ░│",
		"│░░░░░░░░░░░│",
		"└───────────┘",
		"Face down",
	}
}

func cardInteriorRows(rank string, suit string) []string {
	blank := strings.Repeat(" ", 11)
	pair := "  " + suit + "     " + suit + "  "
	centerSuit := centerVisible(suit, 11)
	tripleSuit := centerVisible(suit+suit+suit, 11)

	switch rank {
	case "A":
		return []string{
			padRightVisible(rank, 11),
			blank,
			blank,
			centerSuit,
			blank,
			blank,
			padLeftVisible(rank, 11),
		}
	case "2":
		return []string{
			padRightVisible(rank, 11),
			blank,
			centerSuit,
			blank,
			blank,
			centerSuit,
			padLeftVisible(rank, 11),
		}
	case "3":
		return []string{
			padRightVisible(rank, 11),
			centerSuit,
			blank,
			centerSuit,
			blank,
			centerSuit,
			padLeftVisible(rank, 11),
		}
	case "4":
		return []string{
			padRightVisible(rank, 11),
			pair,
			blank,
			blank,
			blank,
			pair,
			padLeftVisible(rank, 11),
		}
	case "5":
		return []string{
			padRightVisible(rank, 11),
			pair,
			blank,
			centerSuit,
			blank,
			pair,
			padLeftVisible(rank, 11),
		}
	case "6":
		return []string{
			padRightVisible(rank, 11),
			pair,
			blank,
			pair,
			blank,
			pair,
			padLeftVisible(rank, 11),
		}
	case "7":
		return []string{
			padRightVisible(rank, 11),
			pair,
			centerSuit,
			pair,
			blank,
			pair,
			padLeftVisible(rank, 11),
		}
	case "8":
		return []string{
			padRightVisible(rank, 11),
			pair,
			centerSuit,
			pair,
			centerSuit,
			pair,
			padLeftVisible(rank, 11),
		}
	case "9":
		return []string{
			padRightVisible(rank, 11),
			pair,
			pair,
			centerSuit,
			pair,
			pair,
			padLeftVisible(rank, 11),
		}
	case "10":
		return []string{
			padRightVisible(rank, 11),
			pair,
			centerSuit,
			pair,
			pair,
			centerSuit,
			"  " + suit + "     " + suit + rank,
		}
	case "J":
		return []string{
			padRightVisible(rank, 11),
			pair,
			centerVisible("J", 11),
			tripleSuit,
			centerVisible("J", 11),
			pair,
			padLeftVisible(rank, 11),
		}
	case "Q":
		return []string{
			padRightVisible(rank, 11),
			pair,
			centerVisible(".-.", 11),
			centerVisible("( Q )", 11),
			centerVisible("`-'", 11),
			pair,
			padLeftVisible(rank, 11),
		}
	case "K":
		return []string{
			padRightVisible(rank, 11),
			pair,
			centerVisible("\\|/", 11),
			centerVisible("--K--", 11),
			centerVisible("/|\\", 11),
			pair,
			padLeftVisible(rank, 11),
		}
	default:
		return []string{
			padRightVisible(rank, 11),
			blank,
			blank,
			centerSuit,
			blank,
			blank,
			padLeftVisible(rank, 11),
		}
	}
}

func shortRank(rank string) string {
	switch rank {
	case "Ace":
		return "A"
	case "King":
		return "K"
	case "Queen":
		return "Q"
	case "Jack":
		return "J"
	default:
		return rank
	}
}

func suitSymbol(suit string) string {
	switch suit {
	case "Hearts":
		return "♥"
	case "Diamonds":
		return "♦"
	case "Clubs":
		return "♣"
	case "Spades":
		return "♠"
	default:
		return "?"
	}
}

func centerVisible(text string, width int) string {
	padding := max(0, width-lipgloss.Width(text))
	left := padding / 2
	right := padding - left
	return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
}

func padRightVisible(text string, width int) string {
	return text + strings.Repeat(" ", max(0, width-lipgloss.Width(text)))
}

func padLeftVisible(text string, width int) string {
	return strings.Repeat(" ", max(0, width-lipgloss.Width(text))) + text
}
