package main

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vieitesss/gotype/style"
)

type playKeyMaps struct {
	Quit key.Binding
}

func (k playKeyMaps) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k playKeyMaps) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit},
	}
}

var playKeys = playKeyMaps{
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc/C-c", "quit"),
	),
}

// Used when playing, typing.
type PlayHandler struct {
	keys          playKeyMaps
	help          help.Model
	textInput     textinput.Model
	timer         timer.Model
	words         []string
	wordsToRender []string
	currentWord   int
	started       bool
	seconds       int
}

var (
	seconds     = 5
	tiCharLimit = 20
	tiWidth     = 20
)

func NewPlay() *PlayHandler {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = tiCharLimit
	ti.Width = tiWidth

	p := &PlayHandler{
		keys:      playKeys,
		help:      help.New(),
		textInput: ti,
		seconds:   seconds,
		timer:     timer.New(time.Duration(seconds) * time.Second),
	}

	return p
}

func (p PlayHandler) Init() tea.Cmd {
	return tea.Batch(
		getRandomTextToWords,
		textinput.Blink,
	)
}

func (p PlayHandler) Messenger(msg tea.Msg) (Handler, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		p.timer, cmd = p.timer.Update(msg)

		return p, cmd

	case UpdatedWordToRenderMsg:
		p.wordsToRender[p.currentWord] = string(msg)

		return p, nil

	case TextToWriteMsg:
		p.words = msg
		p.wordsToRender = style.InitialWordsStyling(p.words)

		return p, nil

	case tea.KeyMsg:
		// For key maps.
		switch {
		case key.Matches(msg, p.keys.Quit):
			// TODO: handle timer and anything needed,
			//       may ask for confirmation
			return p, updateStatus(Start)
		}

		switch msg.Type {
		case tea.KeySpace:
			p.updateCurrentWord(false)

			// Go to the next word.
			p.currentWord++
			p.textInput.Reset()

			// TODO: show the results.

			// No more words
			if p.currentWord == len(p.words) {
				return p, updateStatus(Quit)
			}

			return p, nil

		default:
			if !p.started {
				p.started = true
				cmds = append(cmds, p.timer.Init())
			}
		}
	}

	if len(p.words) > 0 {
		p.updateCurrentWord(true)
	}

	// Update the input.
	updatedTextInput, cmd := p.textInput.Update(msg)
	p.textInput = updatedTextInput
	cmds = append(cmds, cmd)

	return p, tea.Batch(cmds...)
}

func (p *PlayHandler) updateCurrentWord(addCursor bool) {
	p.wordsToRender[p.currentWord] = style.CompareWithStyle(
		p.textInput.Value(),
		p.words[p.currentWord],
		addCursor,
	)
}

func (p PlayHandler) Render() string {
	s := ""

	// Timer
	if !p.timer.Timedout() {
		s += p.timer.View() + "\n"
	} else {
		s += "Time is over!\n"
	}

	// Words
	if len(p.wordsToRender) > 0 {
		s += strings.Join(p.wordsToRender, " ") + "\n"
	} else {
		s += "the list of words is empty \n"
	}

	// Input
	s += "\n" + p.textInput.View() + "\n"

	// Help
	s += strings.Repeat("\n", 3) + "  " + p.help.View(p.keys)

	return s
}
