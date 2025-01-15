package main

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vieitesss/gotype/style"
)

// Not really stablished yet.
// TextInput defines the text that the user types.
// Words are the words the user has to type.
// WordsToRender are the words the user has to type, but with the corresponding
// color, depending if the user has typed correctly or not.
// CurrentWord indicates the word currently being typed.
// Cmd is the command to execute.
type PlayHandler struct {
	textInput     textinput.Model
	timer         timer.Model
	words         []string
	wordsToRender []string
	currentWord   int
	started       bool
	end           bool
	seconds       int
}

var seconds = 5
var tiCharLimit = 20
var tiWidth = 20

func NewPlay() *PlayHandler {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = tiCharLimit
	ti.Width = tiWidth

	p := &PlayHandler{
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

		if p.timer.Timedout() {
			p.end = true
		}

		return p, cmd

	case UpdatedWordToRenderMsg:
		p.wordsToRender[p.currentWord] = string(msg)
		return p, nil

	case TextToWriteMsg:
		// Save words obtained
		p.words = msg

		// Style the words
		p.wordsToRender = style.InitialWordsStyling(p.words)

		return p, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			// TODO: handle timer and anything needed,
			//       may ask for confirmation
			return p, updateStatus(Start)
		case tea.KeySpace:
			p.updateCurrentWord(false)

			// Go to the next word.
			p.currentWord++
			p.textInput.Reset()

			// TODO: show the results.

			if p.isFinished() {
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

func (p PlayHandler) isFinished() bool {
	return p.currentWord == len(p.words)
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
	s += "\n" + p.textInput.View()

	return s
}
