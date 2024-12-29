package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
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
	words         []string
	wordsToRender []string
	currentWord   int
	typed         string
	cmd           tea.Cmd
}

func NewPlay() *PlayHandler {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 20
	ti.Width = 20

	p := &PlayHandler{
		textInput: ti,
		words:     getRandomTextToWords(),
		cmd:       nil,
	}

	p.styleWords()

	return p
}

// TODO: make this function select an amount of words randomly from a list of
// words.
func getRandomTextToWords() []string {
	return strings.Split("Hola me llamo Dani", " ")
}

// Apply the initial styling to the words.
func (p *PlayHandler) styleWords() {
	p.wordsToRender = make([]string, len(p.words))
	for i, w := range p.words {
		p.wordsToRender[i] = style.Text(w)
	}
}

func (p PlayHandler) Init() tea.Cmd {
	return textinput.Blink
}

func (p *PlayHandler) GetCmd() tea.Cmd {
	cmd := p.cmd
	p.cmd = nil
	return cmd
}

func (p *PlayHandler) Messenger(msg tea.Msg) Action {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return ToQuit
		case tea.KeySpace:
			p.typed = strings.TrimSpace(p.textInput.Value())
			p.currentWord++
			p.textInput.Reset()

			// TODO: show the results.
			if p.isFinished() {
				return ToQuit
			}

			return None
		}
	}

	p.textInput, p.cmd = p.textInput.Update(msg)

	return None
}

func (p *PlayHandler) renderWords() string {
	if p.isFinished() {
		return ""
	}

	if !p.isEmptyInput() {
		p.wordsToRender = p.updatedWords()
	}

	return strings.Join(p.wordsToRender, " ")
}

func (p PlayHandler) isFinished() bool {
	return p.currentWord == len(p.words)
}

func (p PlayHandler) isEmptyInput() bool {
	return len(p.textInput.Value()) == 0
}

// Update the status of the words to render.
func (p *PlayHandler) updatedWords() []string {
	updatedWords := make([]string, len(p.wordsToRender))

	// Add first words.
	for i := 0; i < p.currentWord; i++ {
		updatedWords[i] = p.wordsToRender[i]
	}

	// Update the status of the current word.
	updatedWords[p.currentWord] = p.currentWordStatus()

	// Add the remaining words.
	for i := p.currentWord + 1; i < len(p.words); i++ {
		updatedWords[i] = p.wordsToRender[i]
	}

	return updatedWords
}

// Styles the current word depending on the user input.
func (p PlayHandler) currentWordStatus() string {
	input := p.textInput.Value()
	current := p.words[p.currentWord]
	chars := min(len(input), len(current))
	result := ""

	for i := 0; i < chars; i++ {
		if current[i] == input[i] {
			result += style.Correct(string(current[i]))
		} else {
			result += style.Incorrect(string(current[i]))
		}
	}

	if len(input) < len(current) {
		result += style.Text(current[chars:])
	} else if len(input) > len(current) {
		result += style.Incorrect(input[len(current):])
	}

	return result
}

func (p *PlayHandler) Render() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		p.renderWords(),
		p.textInput.View(),
	) + "\n"
}
