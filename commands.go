package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func updateStatus(status GameStatus) tea.Cmd {
	return func() tea.Msg {
		return status
	}
}

// TODO: make this function select an amount of words randomly from a list of
// words.
func getRandomTextToWords() tea.Msg {
	return textToWrite(strings.Split("Hola me llamo Dani", " "))
}

func quit() tea.Msg {
	return Quit
}