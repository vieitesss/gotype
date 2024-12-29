package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if _, err := tea.NewProgram(NewGame()).Run(); err != nil {
		os.Exit(1)
	}
}
