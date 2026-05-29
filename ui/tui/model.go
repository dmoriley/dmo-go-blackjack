package tui

import (
	"fmt"
	"strconv"
	"strings"

	"blackjack/decks"
	"blackjack/game"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const maxLogEntries = 5

type screen string

const (
	screenNameEntry screen = "name_entry"
	screenTable     screen = "table"
	screenEnd       screen = "end"
)

type logTone string

const (
	logToneInfo    logTone = "info"
	logToneSuccess logTone = "success"
	logToneWarning logTone = "warning"
	logToneDanger  logTone = "danger"
)

type logEntry struct {
	text string
	tone logTone
}

type keyMap struct {
	hit         string
	stand       string
	double      string
	split       string
	continueKey string
	previousBet string
	cashOut     string
	newGame     string
	quit        string
}

type Model struct {
	screen     screen
	engine     *game.Engine
	snapshot   game.Snapshot
	playerName string
	endPhase   game.Phase
	logs       []logEntry
	errMsg     string
	width      int
	height     int
	nameInput  textinput.Model
	betInput   textinput.Model
	keys       keyMap
	styles     styles
}

func NewModel() Model {
	nameInput := textinput.New()
	nameInput.Prompt = "Name: "
	nameInput.Placeholder = "Player"
	nameInput.CharLimit = 24
	nameInput.Focus()

	betInput := textinput.New()
	betInput.Prompt = "$ "
	betInput.Placeholder = "25"
	betInput.CharLimit = 8
	betInput.Width = 12
	betInput.Blur()

	return Model{
		screen:    screenNameEntry,
		width:     120,
		height:    40,
		nameInput: nameInput,
		betInput:  betInput,
		keys: keyMap{
			hit:         "h",
			stand:       "s",
			double:      "d",
			split:       "l",
			continueKey: "enter",
			previousBet: ".",
			cashOut:     "q",
			newGame:     "n",
			quit:        "q",
		},
		styles: newStyles(),
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}

		m.errMsg = ""
		if next, cmd, handled := m.handleKey(msg.String()); handled {
			return next, cmd
		}
	}

	var cmd tea.Cmd
	switch m.screen {
	case screenNameEntry:
		m.nameInput, cmd = m.nameInput.Update(msg)
	case screenTable:
		if m.snapshot.Phase == game.PhaseBetting {
			m.betInput, cmd = m.betInput.Update(msg)
		}
	}

	return m, cmd
}

func (m Model) handleKey(key string) (Model, tea.Cmd, bool) {
	switch m.screen {
	case screenNameEntry:
		switch key {
		case "enter":
			name := strings.TrimSpace(m.nameInput.Value())
			if name == "" {
				m.errMsg = "Enter a player name to start."
				return m, nil, true
			}
			if err := m.startGame(name); err != nil {
				m.errMsg = err.Error()
				return m, nil, true
			}
			m.pushLog(logToneInfo, fmt.Sprintf("Welcome to the table, %s.", name))
			return m, nil, true
		case "esc":
			return m, tea.Quit, true
		}

	case screenTable:
		switch m.snapshot.Phase {
		case game.PhaseBetting:
			switch key {
			case "enter":
				m.placeBetFromInput()
				return m, nil, true
			case "esc":
				m.betInput.SetValue("")
				return m, nil, true
			case m.keys.previousBet:
				if m.snapshot.PreviousBet == 0 {
					m.errMsg = "No previous bet is available yet."
					return m, nil, true
				}
				m.placeBet(m.snapshot.PreviousBet)
				return m, nil, true
			case m.keys.cashOut:
				m.cashOut()
				return m, nil, true
			}

		case game.PhasePlayerTurn:
			switch key {
			case m.keys.hit:
				m.applyMove(game.MoveHit, "Hit is not available right now.")
				return m, nil, true
			case m.keys.stand:
				m.applyMove(game.MoveStand, "Stand is not available right now.")
				return m, nil, true
			case m.keys.double:
				m.applyMove(game.MoveDouble, "Double is not available for this hand.")
				return m, nil, true
			case m.keys.split:
				m.applyMove(game.MoveSplit, "Split is not available for this hand.")
				return m, nil, true
			}

		case game.PhaseRoundResult:
			if key == m.keys.continueKey {
				m.continueRound()
				return m, nil, true
			}
		}

	case screenEnd:
		switch key {
		case m.keys.newGame:
			if err := m.startGame(m.playerName); err != nil {
				m.errMsg = err.Error()
				return m, nil, true
			}
			m.pushLog(logToneInfo, fmt.Sprintf("A new shoe is ready for %s.", m.playerName))
			return m, nil, true
		case m.keys.quit, "esc":
			return m, tea.Quit, true
		}
	}

	return m, nil, false
}

func (m *Model) startGame(playerName string) error {
	engine, err := game.NewEngine(game.Config{
		PlayerName:   playerName,
		StartingCash: 500,
		DeckConfig:   decks.NewBlackjackDeckConfig().WithNumberOfDecks(4),
		ShuffleCount: 5,
	})
	if err != nil {
		return err
	}

	m.engine = engine
	m.playerName = playerName
	m.snapshot = engine.Snapshot()
	m.screen = screenTable
	m.endPhase = ""
	m.logs = nil
	m.errMsg = ""
	m.betInput.SetValue("")
	m.nameInput.SetValue(playerName)
	m.syncInputs()
	return nil
}

func (m *Model) syncInputs() {
	if m.screen == screenNameEntry {
		m.nameInput.Focus()
		m.betInput.Blur()
		return
	}

	m.nameInput.Blur()
	if m.screen == screenTable && m.snapshot.Phase == game.PhaseBetting {
		m.betInput.Focus()
		return
	}

	m.betInput.Blur()
}

func (m *Model) placeBetFromInput() {
	amount, err := strconv.Atoi(strings.TrimSpace(m.betInput.Value()))
	if err != nil {
		m.errMsg = "Enter a valid numeric bet."
		return
	}
	m.placeBet(amount)
}

func (m *Model) placeBet(amount int) {
	if m.engine == nil {
		m.errMsg = "No active game is running."
		return
	}

	result, err := m.engine.PlaceBet(amount)
	if err != nil {
		m.errMsg = err.Error()
		return
	}

	m.betInput.SetValue("")
	m.applyResult(result)
}

func (m *Model) applyMove(move game.Move, illegalMessage string) {
	if !m.hasLegalMove(move) {
		m.errMsg = illegalMessage
		return
	}

	result, err := m.engine.ApplyMove(move)
	if err != nil {
		m.errMsg = err.Error()
		return
	}

	m.applyResult(result)
}

func (m *Model) continueRound() {
	result, err := m.engine.Continue()
	if err != nil {
		m.errMsg = err.Error()
		return
	}

	m.applyResult(result)
}

func (m *Model) cashOut() {
	result, err := m.engine.CashOut()
	if err != nil {
		m.errMsg = err.Error()
		return
	}

	m.applyResult(result)
}

func (m *Model) applyResult(result game.StepResult) {
	m.snapshot = result.Snapshot
	m.addEvents(result.Events)

	switch result.Snapshot.Phase {
	case game.PhaseGameOver, game.PhaseCashedOut:
		m.screen = screenEnd
		m.endPhase = result.Snapshot.Phase
	default:
		m.screen = screenTable
		m.endPhase = ""
	}

	m.syncInputs()
}

func (m *Model) hasLegalMove(target game.Move) bool {
	for _, move := range m.snapshot.LegalMoves {
		if move == target {
			return true
		}
	}
	return false
}

func (m *Model) addEvents(events []game.Event) {
	for _, event := range events {
		for _, entry := range translateEvent(event, len(m.snapshot.PlayerHands)) {
			m.pushLog(entry.tone, entry.text)
		}
	}
}

func (m *Model) pushLog(tone logTone, text string) {
	cleaned := strings.TrimSpace(strings.Join(strings.Fields(text), " "))
	if cleaned == "" {
		return
	}

	m.logs = append(m.logs, logEntry{tone: tone, text: cleaned})
	if len(m.logs) > maxLogEntries {
		m.logs = append([]logEntry(nil), m.logs[len(m.logs)-maxLogEntries:]...)
	}
}

func translateEvent(event game.Event, handCount int) []logEntry {
	switch e := event.(type) {
	case game.ReshuffledEvent:
		return []logEntry{{tone: logToneWarning, text: fmt.Sprintf("The shoe was reshuffled. %d cards remain.", e.Remaining)}}
	case game.HandFinishedEvent:
		hand := formatHandLabel(e.HandIndex, handCount)
		switch e.Reason {
		case game.HandFinishBust:
			return []logEntry{{tone: logToneDanger, text: fmt.Sprintf("%s busted at %d.", hand, e.Total)}}
		case game.HandFinishTwentyOne:
			return []logEntry{{tone: logToneSuccess, text: fmt.Sprintf("%s reached 21.", hand)}}
		case game.HandFinishSplitAceAutoStand:
			return []logEntry{{tone: logToneInfo, text: fmt.Sprintf("%s was dealt as split aces and stands automatically.", hand)}}
		}
	case game.HandResolvedEvent:
		hand := formatHandLabel(e.HandIndex, handCount)
		switch e.Outcome {
		case game.OutcomeWon:
			if e.Blackjack {
				return []logEntry{{tone: logToneSuccess, text: fmt.Sprintf("%s won $%d against the dealer %d to %d with blackjack.", hand, e.Payout, e.PlayerTotal, e.DealerTotal)}}
			}
			return []logEntry{{tone: logToneSuccess, text: fmt.Sprintf("%s won $%d against the dealer %d to %d.", hand, e.Payout, e.PlayerTotal, e.DealerTotal)}}
		case game.OutcomeLost:
			return []logEntry{{tone: logToneDanger, text: fmt.Sprintf("%s lost $%d to the dealer %d to %d.", hand, e.Bet, e.PlayerTotal, e.DealerTotal)}}
		case game.OutcomePush:
			return []logEntry{{tone: logToneInfo, text: fmt.Sprintf("%s pushed and returned $%d.", hand, e.Payout)}}
		}
	case game.ActionDeniedEvent:
		return []logEntry{{tone: logToneWarning, text: e.Reason}}
	case game.CashOutEvent:
		return []logEntry{{tone: logToneSuccess, text: fmt.Sprintf("%s cashes out with $%d.", e.PlayerName, e.Cash)}}
	case game.GameOverEvent:
		return []logEntry{{tone: logToneDanger, text: fmt.Sprintf("%s is out of cash.", e.PlayerName)}}
	}

	return nil
}

func formatHandLabel(index int, handCount int) string {
	if handCount <= 1 {
		return "Your hand"
	}
	return fmt.Sprintf("Hand %d", index+1)
}
