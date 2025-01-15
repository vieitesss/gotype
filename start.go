package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vieitesss/gotype/style"
)

// The main menu.
// For now you can choose between playing or exiting de app.
// Current indicates currently selected option.
// Choices are the posible options in the menu.
// Actions are the actions corresponding to the choices.
type StartHandler struct {
	current  int
	choices  []string
	statuses []GameStatus
}

func NewStart() *StartHandler {
	s := &StartHandler{
		choices:  []string{"Play", "Quit"},
		statuses: []GameStatus{Playing, Quit},
	}

	return s
}

func (s StartHandler) Init() tea.Cmd {
	return nil
}

func (s StartHandler) Messenger(msg tea.Msg) (Handler, tea.Cmd) {
	n := len(s.choices)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDown:
			s.current = (s.current + 1) % n
		case tea.KeyUp:
			s.current -= 1
			if s.current < 0 {
				s.current = n - 1
			}

		case tea.KeyEnter:
			status := s.statuses[s.current]

			return s, updateStatus(status)

		case tea.KeyCtrlC:
			return s, updateStatus(Quit)
		}
	}

	return s, nil
}

func (s StartHandler) Render() string {
	res := ""

	for i, choice := range s.choices {
		if i == s.current {
			res += style.MainMenuOptionStyling(" " + choice)
		} else {
			res += style.MainMenuOptionStyling(choice)
		}
		res += "\n"
	}

	return res
}
