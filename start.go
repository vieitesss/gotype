package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vieitesss/gotype/style"
)

// The main menu.
// For now you can choose between playing or exiting de app.
// Current indicates currently selected option.
// Choices are the posible options in the menu.
// Actions are the actions corresponding to the choices.
type StartHandler struct {
	current int
	choices []string
	actions []Action
}

func NewStart() *StartHandler {
	s := &StartHandler{
		actions: []Action{ToPlaying, ToQuit},
		choices: []string{"Play", "Quit"},
	}

	return s
}

func (s StartHandler) Init() tea.Cmd {
	return nil
}

func (s StartHandler) GetCmd() tea.Cmd {
	return nil
}

func (s *StartHandler) Messenger(msg tea.Msg) Action {
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
			return s.actions[s.current]
		case tea.KeyCtrlC:
			return ToQuit
		}
	}
	return None
}

func (s StartHandler) Render() string {
	return fmt.Sprintf("%s", s.printList()) + "\n"
}

func (s StartHandler) printList() string {
	result := ""

	for i, c := range s.choices {
		selector := "  "
		if i == s.current {
			selector = " "
		}

		result += fmt.Sprintf("%s%s\n", style.Selector(selector), style.Text(c))
	}

	return result
}
