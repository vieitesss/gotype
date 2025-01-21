package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vieitesss/gotype/style"
)

const (
	listHeight   = 8
	defaultWidth = 20
)

type item string

func (i item) FilterValue() string { return string(i) }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	styledItem := style.MainMenuOptionStyling(string(i), index, m.Index())

	fmt.Fprint(w, styledItem)
}

// The main menu.
type StartHandler struct {
	menu list.Model
}

func NewStart() *StartHandler {
	options := []list.Item{
		item("Play"),
		item("Quit"),
	}

	menu := list.New(options, itemDelegate{}, defaultWidth, listHeight)
	menu.Title = "Main menu"
	menu.Styles.Title = style.MenuStyle
	menu.SetShowStatusBar(false)
	menu.SetFilteringEnabled(false)

	s := &StartHandler{
		menu: menu,
	}

	return s
}

func (s StartHandler) Init() tea.Cmd {
	return nil
}

func (s StartHandler) Messenger(msg tea.Msg) (Handler, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyEnter:
			var status GameStatus
			switch s.menu.SelectedItem().FilterValue() {

			case "Play":
				status = Playing

			case "Quit":
				status = Quit
			}

			return s, updateStatus(status)

		case tea.KeyCtrlC:
			return s, updateStatus(Quit)
		}
	}

	var cmd tea.Cmd
	s.menu, cmd = s.menu.Update(msg)
	return s, cmd
}

func (s StartHandler) Render() string {
	return "\n" + s.menu.View()
}
