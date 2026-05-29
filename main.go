package main

import (
	"blackjack/ui/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	program := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
