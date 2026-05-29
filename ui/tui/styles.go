package tui

import "github.com/charmbracelet/lipgloss"

type styles struct {
	app             lipgloss.Style
	canvas          lipgloss.Style
	shell           lipgloss.Style
	header          lipgloss.Style
	headerTitle     lipgloss.Style
	headerMeta      lipgloss.Style
	panel           lipgloss.Style
	panelTitle      lipgloss.Style
	muted           lipgloss.Style
	accent          lipgloss.Style
	success         lipgloss.Style
	warning         lipgloss.Style
	danger          lipgloss.Style
	input           lipgloss.Style
	inputLabel      lipgloss.Style
	card            lipgloss.Style
	cardFaceDown    lipgloss.Style
	cardAccent      lipgloss.Style
	activeHand      lipgloss.Style
	inactiveHand    lipgloss.Style
	actionBar       lipgloss.Style
	actionKey       lipgloss.Style
	actionDisabled  lipgloss.Style
	logPanel        lipgloss.Style
	endPanel        lipgloss.Style
	namePrompt      lipgloss.Style
	nameShell       lipgloss.Style
	divider         lipgloss.Style
	statusPill      lipgloss.Style
	statusPillAlert lipgloss.Style
}

func newStyles() styles {
	base := lipgloss.Color("#E9E4D8")
	muted := lipgloss.Color("#9EA79C")
	green := lipgloss.Color("#153C34")
	greenAlt := lipgloss.Color("#1E5448")
	gold := lipgloss.Color("#D7B46A")
	amber := lipgloss.Color("#E2A65A")
	red := lipgloss.Color("#C06B65")
	border := lipgloss.Color("#355B52")

	return styles{
		app: lipgloss.NewStyle().
			Foreground(base).
			Background(green).
			Padding(1, 2),
		canvas: lipgloss.NewStyle().
			Foreground(base).
			Background(green),
		shell: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(border).
			Padding(1, 2),
		header: lipgloss.NewStyle().
			MarginBottom(1),
		headerTitle: lipgloss.NewStyle().
			Foreground(gold).
			Bold(true),
		headerMeta: lipgloss.NewStyle().
			Foreground(muted),
		panel: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(border).
			Padding(1, 2),
		panelTitle: lipgloss.NewStyle().
			Foreground(gold).
			Bold(true),
		muted:   lipgloss.NewStyle().Foreground(muted),
		accent:  lipgloss.NewStyle().Foreground(gold).Bold(true),
		success: lipgloss.NewStyle().Foreground(lipgloss.Color("#9ED39A")),
		warning: lipgloss.NewStyle().Foreground(amber),
		danger:  lipgloss.NewStyle().Foreground(red),
		input: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(gold).
			Padding(0, 1),
		inputLabel: lipgloss.NewStyle().Foreground(muted),
		card: lipgloss.NewStyle().
			Width(18).
			Height(5).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#6A8A7F")).
			Padding(0, 1),
		cardFaceDown: lipgloss.NewStyle().
			Width(18).
			Height(5).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4D675F")).
			Foreground(muted).
			Padding(0, 1),
		cardAccent: lipgloss.NewStyle().Foreground(gold).Bold(true),
		activeHand: lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(gold).
			Padding(1, 2),
		inactiveHand: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(border).
			Padding(1, 2),
		actionBar: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(border).
			Padding(1, 2),
		actionKey:      lipgloss.NewStyle().Foreground(gold).Bold(true),
		actionDisabled: lipgloss.NewStyle().Foreground(muted),
		logPanel: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(border).
			Padding(1, 2),
		endPanel: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(gold).
			Padding(1, 3),
		namePrompt: lipgloss.NewStyle().
			Foreground(gold).
			Bold(true),
		nameShell: lipgloss.NewStyle().
			Width(56).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(border).
			Padding(2, 3),
		divider: lipgloss.NewStyle().Foreground(border),
		statusPill: lipgloss.NewStyle().
			Foreground(greenAlt).
			Background(gold).
			Padding(0, 1),
		statusPillAlert: lipgloss.NewStyle().
			Foreground(greenAlt).
			Background(amber).
			Padding(0, 1),
	}
}
