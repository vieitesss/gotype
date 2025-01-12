package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Each Handler belongs to a specific status of the game.
// Does the necessary stuff while in that status.
type Handler interface {
	Messenger(msg tea.Msg) (Handler, tea.Cmd)
	Render() string
	Init() tea.Cmd
}
