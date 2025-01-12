package main

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vieitesss/gotype/style"
)

type textToWrite []string
type updatedWordToRender string

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

	case updatedWordToRender:
		p.wordsToRender[p.currentWord] = string(msg)
		return p, nil

	case textToWrite:
		// Save words obtained
		p.words = msg

		// Style the words
		p.wordsToRender = make([]string, len(p.words))
		for i, w := range p.words {
			p.wordsToRender[i] = style.Text(w)
		}
		return p, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			// TODO: handle timer and anything needed,
			//       may ask for confirmation
			return p, updateStatus(Start)
		case tea.KeySpace:
			// Update the current word to remove the cursor.
			p.wordsToRender[p.currentWord] = p.updateWord(false)

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
		// Update the current word status whit the current input after a key press.
		p.wordsToRender[p.currentWord] = p.updateWord(true)
	}

	// Update the input.
	updatedTextInput, cmd := p.textInput.Update(msg)
	p.textInput = updatedTextInput
	cmds = append(cmds, cmd)

	return p, tea.Batch(cmds...)
}

func (p PlayHandler) updateWord(addCursor bool) string {
	input := p.textInput.Value()
	word := p.words[p.currentWord]

	toRender := p.textStatus(input, p.currentWord)

	if len(input) >= len(word) {
		return toRender
	}

	if addCursor {
		toRender += style.CurrentChar(string(word[len(input)]))
		if len(word) > len(input)+1 {
			toRender += style.Text(word[len(input)+1:])
		}
	} else {
		toRender += style.Text(word[len(input):])
	}

	return toRender
}

// Gets a text to be checked and the index of the word to be checked with.
// Returns the text styled accordingly.
func (p PlayHandler) textStatus(text string, wordIndex int) string {
	current := p.words[wordIndex]
	chars := min(len(text), len(current))
	result := ""

	// Print already written chars
	for i := 0; i < chars; i++ {
		if text[i] == current[i] {
			result += style.Correct(string(current[i]))
		} else {
			result += style.Incorrect(string(current[i]))
		}
	}

	// Print the excess chars, if there are any.
	if len(text) > len(current) {
		result += style.Incorrect(text[len(current):])
	}

	return result
}

func (p PlayHandler) isFinished() bool {
	return p.currentWord == len(p.words)
}

func (p PlayHandler) Render() string {
	s := ""

	if !p.timer.Timedout() {
		s += p.timer.View() + "\n"
	} else {
		s += "Time is over!\n"
	}

	if len(p.wordsToRender) > 0 {
		s += strings.Join(p.wordsToRender, " ") + "\n"
	} else {
		s += "the list is empty"
	}

	s += "\n" + p.textInput.View()

	return s
}
