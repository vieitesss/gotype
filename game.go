package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// The current status of the game.
type GameStatus int

const (
	// Initial menu.
	Start GameStatus = iota

	// While the user is playing, typing.
	Playing

	// When ending the game.
	Quit
)

// Main model
type Game struct {
	Status   GameStatus
	Handlers map[GameStatus]Handler
}

func NewGame() Game {
	g := Game{Status: Start}

	l := make(map[GameStatus]Handler)
	l[Start] = NewStart()
	l[Playing] = NewPlay()

	g.Handlers = l

	return g
}

func (g Game) Init() tea.Cmd {
	return nil
}

func (g Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case GameStatusUpdatedMsg:
		switch msg.status {
		case Quit:
			return g, tea.Quit

		default:
			g.Status = msg.status

			return g, g.Handlers[msg.status].Init()
		}
	}

	updated, cmd := g.Handlers[g.Status].Messenger(msg)
	g.Handlers[g.Status] = updated

	return g, cmd
}

func (g Game) View() string {
	return g.Handlers[g.Status].Render()
}
