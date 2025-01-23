package keymaps

import (
	"github.com/charmbracelet/bubbles/key"
)

type PlayKeyMaps struct {
	Quit key.Binding
}

func (k PlayKeyMaps) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k PlayKeyMaps) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit},
	}
}

func PlayKeys() PlayKeyMaps {
	return PlayKeyMaps{
		Quit: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc/C-c", "quit"),
		),
	}
}
