package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Each Handler belongs to a specific status of the game.
// Does the necessary stuff while in that status.
type Handler interface {
	Messenger(tea.Msg) Action
	Render() string
	Init() tea.Cmd
	GetCmd() tea.Cmd
}

// The current status of the game.
type GameStatus int

const (
	// Initial menu.
	Start GameStatus = iota

	// While the user is plaing, typing.
	Playing

	// When ending the game.
	Quit
)

// The action to make.
type Action int

const (
	// Changing to "Start" status.
	ToStart Action = iota

	// Changing to "Play" status.
	ToPlaying

	// Changing to "Quit" status.
	ToQuit

	// When no action is required, this is returned.
	None
)

// The game definition.
// Status defines the current game status.
// The handlers are a map linking each status with the handler that corresponds
// to it.
type Game struct {
	Status   GameStatus
	Handlers map[GameStatus]Handler
}

// Initializes the game.
func NewGame() Game {
	g := Game{Status: Start}

	l := make(map[GameStatus]Handler)
	l[Start] = NewStart()
	l[Playing] = NewPlay()

	g.Handlers = l

	return g
}

func (g Game) Init() tea.Cmd {
	return g.Handlers[Playing].Init()
}

func (g Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	action := g.Handlers[g.Status].Messenger(msg)
	currentStatus := g.Status
	cmd := g.Handlers[currentStatus].GetCmd()

	switch action {
	case ToStart:
		g.Status = Start
	case ToPlaying:
		g.Status = Playing
	case ToQuit:
		return g, tea.Quit
	}

	return g, cmd
}

func (g Game) View() string {
	return g.Handlers[g.Status].Render()
}
