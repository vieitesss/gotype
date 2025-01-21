package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	program := tea.NewProgram(
		NewGame(),
		tea.WithAltScreen(),
	)

	if _, err := program.Run(); err != nil {
		os.Exit(1)
	}
}
