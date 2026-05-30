package tui

import (
	"fmt"
	"strings"

	"blackjack/game"

	"github.com/charmbracelet/lipgloss"
)

type cardRenderMode string

const (
	cardRenderModeWide    cardRenderMode = "wide"
	cardRenderModeCompact cardRenderMode = "compact"
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
	availableWidth, availableHeight := m.appViewport()
	sections := []string{
		m.styles.headerTitle.Render("Blackjack"),
		m.styles.muted.Render("A restrained table view built on the headless engine."),
		"",
		m.styles.namePrompt.Render("Choose the name that will sit at the table."),
		m.styles.input.Render(m.nameInput.View()),
		m.renderError(),
		m.styles.muted.Render("Press enter to start. Press esc to quit."),
	}

	nameShell := fitStyleWidth(m.styles.nameShell.Copy(), min(56, availableWidth))
	card := nameShell.Render(strings.Join(filterEmpty(sections), "\n"))
	placed := lipgloss.Place(
		availableWidth,
		availableHeight,
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
	contentWidth := m.shellContentWidth()
	header := m.viewHeader(contentWidth)
	cardMode := cardRenderModeWide
	if m.shouldStackTable(contentWidth) {
		cardMode = cardRenderModeCompact
	}

	sections := []string{header}
	var body string
	if m.shouldStackTable(contentWidth) {
		table := m.viewMainTable(contentWidth, cardMode)
		sidebar := m.viewSidebar(contentWidth)
		footer := m.viewFooter(contentWidth)
		body = lipgloss.JoinVertical(lipgloss.Left, table, sidebar, footer)
	} else {
		sidebarWidth := m.sidebarWidth(contentWidth)
		tableWidth := max(1, contentWidth-sidebarWidth-2)
		table := m.viewMainTable(tableWidth, cardMode)
		footer := m.viewFooter(tableWidth)
		sidebar := m.viewSidebar(sidebarWidth)
		leftColumn := lipgloss.JoinVertical(lipgloss.Left, table, footer)
		body = lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, "  ", sidebar)
	}

	sections = append(sections, body)
	return fitStyleWidth(m.styles.shell.Copy(), contentWidth+m.styles.shell.GetHorizontalFrameSize()).Render(strings.Join(filterEmpty(sections), "\n"))
}

func (m Model) viewHeader(contentWidth int) string {
	left := m.styles.headerTitle.Render("Blackjack")
	return m.styles.header.Render(left)
}

func (m Model) viewMainTable(width int, cardMode cardRenderMode) string {
	dealer := m.renderHandPanel("Dealer", m.snapshot.Dealer, false, false, 0, width, cardMode)
	hands := []string{dealer}

	for idx, hand := range m.snapshot.PlayerHands {
		label := m.snapshot.PlayerName
		if m.snapshot.IsSplitRound {
			label = fmt.Sprintf("%s · Hand %d", m.snapshot.PlayerName, idx+1)
		}
		active := idx == m.snapshot.ActiveHandIndex && m.snapshot.Phase == game.PhasePlayerTurn
		hands = append(hands, m.renderHandPanel(label, hand, true, active, idx, width, cardMode))
	}

	return lipgloss.JoinVertical(lipgloss.Left, hands...)
}

func (m Model) renderHandPanel(title string, hand game.HandState, isPlayer bool, active bool, idx int, width int, cardMode cardRenderMode) string {
	style := m.styles.inactiveHand
	if active {
		style = m.styles.activeHand
	}
	if width > 0 {
		style = fitStyleWidth(style.Copy(), width)
	}

	innerWidth := max(1, width-style.GetHorizontalFrameSize())
	if width > 0 {
		innerWidth = max(1, width-style.GetHorizontalFrameSize())
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

	titleText := m.styles.panelTitle.Render(title)
	metaText := m.styles.muted.Render(strings.Join(meta, " · "))
	header := titleText
	if metaText != "" {
		if innerWidth > 0 && lipgloss.Width(titleText)+2+lipgloss.Width(metaText) <= innerWidth {
			header = lipgloss.JoinHorizontal(lipgloss.Top, titleText, "  ", metaText)
		} else {
			header = lipgloss.JoinVertical(lipgloss.Left, titleText, metaText)
		}
	}

	return style.Render(lipgloss.JoinVertical(lipgloss.Left,
		header,
		m.renderCards(hand.Cards, innerWidth, cardMode),
	))
}

func (m Model) renderCards(cards []game.CardState, availableWidth int, cardMode cardRenderMode) string {
	if len(cards) == 0 {
		return m.styles.muted.Render("No cards in play")
	}

	rendered := make([]string, 0, len(cards))
	for _, card := range cards {
		rendered = append(rendered, m.renderCard(card, cardMode))
	}

	return wrapBlocks(rendered, availableWidth)
}

func (m Model) renderCard(card game.CardState, cardMode cardRenderMode) string {
	switch cardMode {
	case cardRenderModeCompact:
		if !card.FaceUp {
			return m.styles.cardFaceDown.Copy().Width(12).Render(strings.Join(faceDownCompactCardLines(), "\n"))
		}
		return m.styles.card.Copy().Width(12).Render(strings.Join(faceUpCompactCardLines(card), "\n"))
	default:
		if !card.FaceUp {
			return m.styles.cardFaceDown.Render(strings.Join(faceDownCardLines(), "\n"))
		}

		return m.styles.card.Render(strings.Join(faceUpCardLines(card), "\n"))
	}
}

func (m Model) viewSidebar(width int) string {
	statusStyle := m.styles.panel
	logStyle := m.styles.logPanel
	phaseStyle := m.styles.panel
	if width > 0 {
		statusStyle = fitStyleWidth(statusStyle.Copy(), width)
		logStyle = fitStyleWidth(logStyle.Copy(), width)
		phaseStyle = fitStyleWidth(phaseStyle.Copy(), width)
	}

	phase := phaseLabel(m.snapshot.Phase)
	pill := m.styles.statusPill.Render(phase)
	if m.snapshot.Phase == game.PhaseRoundResult || m.snapshot.Phase == game.PhaseGameOver || m.snapshot.Phase == game.PhaseCashedOut {
		pill = m.styles.statusPillAlert.Render(phase)
	}

	phasePanel := phaseStyle.Render(strings.Join([]string{
		m.styles.panelTitle.Render("Game Status"),
		pill,
		fmt.Sprintf("Deck %d / %d", m.snapshot.DeckRemaining, m.snapshot.DeckTotal),
	}, "\n"))

	statusLines := []string{
		m.styles.panelTitle.Render("Player Status"),
		"Player ----------- " + padLeftVisible(m.playerName, 5),
		"Cash ------------- " + padLeftVisible(fmt.Sprintf("$%d", m.snapshot.Cash), 5),
		"Previous bet ----- " + padLeftVisible(fmt.Sprintf("$%d", m.snapshot.PreviousBet), 5),
		"Pending results -- " + padLeftVisible(fmt.Sprintf("%d", m.snapshot.PendingResults), 5),
	}
	status := statusStyle.Render(strings.Join(statusLines, "\n"))

	logLines := make([]string, 0, len(m.logs)+1)
	logLines = append(logLines, m.styles.panelTitle.Render("Table Log"))
	if len(m.logs) == 0 {
		logLines = append(logLines, m.styles.muted.Render("The dealer is waiting for the next action."))
	} else {
		for _, entry := range m.logs {
			logLines = append(logLines, m.renderLogEntry(entry))
		}
	}

	logPanel := logStyle.Render(strings.Join(logLines, "\n"))
	return lipgloss.JoinVertical(lipgloss.Left, phasePanel, status, logPanel)
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

func (m Model) viewFooter(width int) string {
	barStyle := m.styles.actionBar
	innerWidth := width
	if width > 0 {
		barStyle = fitStyleWidth(barStyle.Copy(), width)
		innerWidth = max(1, width-barStyle.GetHorizontalFrameSize())
	}

	actions := m.renderActions(innerWidth)
	help := []string{actions}
	if msg := m.renderError(); msg != "" {
		help = append(help, msg)
	}
	return barStyle.Render(strings.Join(help, "\n"))
}

func (m Model) renderActions(width int) string {
	switch m.snapshot.Phase {
	case game.PhaseBetting:
		input := lipgloss.JoinHorizontal(lipgloss.Center,
			m.styles.inputLabel.Render("Bet"),
			" ",
			m.styles.input.Render(m.betInput.View()),
		)
		if lipgloss.Width(input) > width {
			input = lipgloss.JoinVertical(lipgloss.Left,
				m.styles.inputLabel.Render("Bet"),
				m.styles.input.Render(m.betInput.View()),
			)
		}
		actions := []string{
			m.renderAction("enter", "place bet", true),
			m.renderAction(m.keys.previousBet, "repeat previous bet", m.snapshot.PreviousBet > 0),
			m.renderAction(m.keys.cashOut, "cash out", true),
		}
		return lipgloss.JoinVertical(lipgloss.Left,
			input,
			wrapBlocks(actions, width),
		)
	case game.PhasePlayerTurn:
		actions := []string{
			m.renderAction(m.keys.hit, "hit", m.hasLegalMove(game.MoveHit)),
			m.renderAction(m.keys.stand, "stand", m.hasLegalMove(game.MoveStand)),
			m.renderAction(m.keys.double, "double", m.hasLegalMove(game.MoveDouble)),
			m.renderAction(m.keys.split, "split", m.hasLegalMove(game.MoveSplit)),
		}
		return wrapBlocks(actions, width)
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
	availableWidth, availableHeight := m.appViewport()
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

	endPanel := fitStyleWidth(m.styles.endPanel.Copy(), min(72, availableWidth))
	panel := endPanel.Render(strings.Join(body, "\n"))
	placed := lipgloss.Place(
		availableWidth,
		availableHeight,
		lipgloss.Center,
		lipgloss.Center,
		panel,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(m.styles.canvas.GetForeground()),
		lipgloss.WithWhitespaceBackground(m.styles.canvas.GetBackground()),
	)
	return m.styles.canvas.Render(placed)
}

func (m Model) appViewport() (int, int) {
	return max(1, m.width-m.styles.app.GetHorizontalFrameSize()), max(1, m.height-m.styles.app.GetVerticalFrameSize())
}

func (m Model) shellContentWidth() int {
	availableWidth, _ := m.appViewport()
	return max(1, availableWidth-m.styles.shell.GetHorizontalFrameSize())
}

func (m Model) shouldStackTable(contentWidth int) bool {
	return contentWidth < 118
}

func (m Model) sidebarWidth(contentWidth int) int {
	return min(38, max(30, contentWidth/3))
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

func wrapBlocks(blocks []string, availableWidth int) string {
	if len(blocks) == 0 {
		return ""
	}
	if availableWidth <= 0 {
		return lipgloss.JoinVertical(lipgloss.Left, blocks...)
	}

	rows := make([]string, 0, len(blocks))
	current := make([]string, 0, len(blocks))
	currentWidth := 0

	for _, block := range blocks {
		blockWidth := lipgloss.Width(block)
		if len(current) > 0 && currentWidth+blockWidth > availableWidth {
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, current...))
			current = nil
			currentWidth = 0
		}

		current = append(current, block)
		currentWidth += blockWidth
	}

	if len(current) > 0 {
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, current...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func fitStyleWidth(style lipgloss.Style, totalWidth int) lipgloss.Style {
	if totalWidth <= 0 {
		return style
	}
	return style.Width(max(1, totalWidth-style.GetHorizontalFrameSize()))
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

func faceUpCompactCardLines(card game.CardState) []string {
	rank := shortRank(card.Rank)
	suit := suitSymbol(card.Suit)
	rows := compactCardInteriorRows(rank, suit)
	lines := make([]string, 0, len(rows)+3)
	lines = append(lines, "┌─────────┐")
	for _, row := range rows {
		lines = append(lines, "│"+row+"│")
	}
	lines = append(lines, "└─────────┘")
	lines = append(lines, compactCaption(card))
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

func faceDownCompactCardLines() []string {
	return []string{
		"┌─────────┐",
		"│░░░░░░░░░│",
		"│░▒▒▒▒▒▒▒░│",
		"│░░░░░░░░░│",
		"└─────────┘",
		"Down",
	}
}

func compactCardInteriorRows(rank string, suit string) []string {
	switch rank {
	case "A":
		return []string{
			"A       A",
			"         ",
			"    " + suit + "    ",
		}
	case "2":
		return []string{
			"2 " + suit + "   " + suit + " 2",
			"         ",
			"         ",
		}
	case "3":
		return []string{
			"3 " + suit + "   " + suit + " 3",
			"    " + suit + "    ",
			"         ",
		}
	case "4":
		return []string{
			"4 " + suit + "   " + suit + " 4",
			"         ",
			"4 " + suit + "   " + suit + " 4",
		}
	case "5":
		return []string{
			"5 " + suit + "   " + suit + " 5",
			"    " + suit + "    ",
			"5 " + suit + "   " + suit + " 5",
		}
	case "6":
		return []string{
			"6 " + suit + "   " + suit + " 6",
			"6 " + suit + "   " + suit + " 6",
			"6 " + suit + "   " + suit + " 6",
		}
	case "7":
		return []string{
			"7 " + suit + "   " + suit + " 7",
			"   " + suit + " " + suit + "   ",
			"7 " + suit + "   " + suit + " 7",
		}
	case "8":
		return []string{
			"8 " + suit + "   " + suit + " 8",
			"8 " + suit + "   " + suit + " 8",
			"8 " + suit + "   " + suit + " 8",
		}
	case "9":
		return []string{
			"9 " + suit + " " + suit + " " + suit + " 9",
			"   " + suit + " " + suit + "   ",
			"9 " + suit + " " + suit + " " + suit + " 9",
		}
	case "10":
		return []string{
			"10" + suit + " " + suit + " " + suit + "10",
			"  " + suit + " " + suit + " " + suit + "  ",
			"10" + suit + " " + suit + " " + suit + "10",
		}
	case "J":
		return []string{
			"J " + suit + "   " + suit + " J",
			"   " + suit + "J" + suit + "   ",
			"J " + suit + "   " + suit + " J",
		}
	case "Q":
		return []string{
			"Q " + suit + "   " + suit + " Q",
			"   (Q)   ",
			"Q " + suit + "   " + suit + " Q",
		}
	case "K":
		return []string{
			"K " + suit + "   " + suit + " K",
			"   <K>   ",
			"K " + suit + "   " + suit + " K",
		}
	default:
		return []string{
			padRightVisible(rank, 9),
			"    " + suit + "    ",
			padLeftVisible(rank, 9),
		}
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

func suitInitial(suit string) string {
	switch suit {
	case "Hearts":
		return "H"
	case "Diamonds":
		return "D"
	case "Clubs":
		return "C"
	case "Spades":
		return "S"
	default:
		return "?"
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

func compactCaption(card game.CardState) string {
	return fmt.Sprintf("%s + %s [%d]", shortRank(card.Rank), suitInitial(card.Suit), card.Value)
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
