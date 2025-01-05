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
			// TODO: handle timer and anything needed,
			//       may ask for confirmation
			return ToStart
		case tea.KeySpace:
			// Update the current word to remove the cursor.
			input := p.textInput.Value()
			word := p.words[p.currentWord]
			toRender := p.textStatus(input)
			if len(input) < len(word) {
				toRender += style.Text(word[len(input):])
			}
			p.wordsToRender[p.currentWord] = toRender

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

// Returns a string joining all the words to be written.
func (p *PlayHandler) renderWords() string {
	if p.isFinished() {
		return ""
	}

	// Get the status of the input.
	input := p.textInput.Value()
	word := p.words[p.currentWord]
	toRender := p.textStatus(input)

	// Add cursor and remaining chars if needed.
	if len(input) < len(word) {
		toRender += style.CurrentChar(string(word[len(input)]))
		if len(input) < len(word)+1 {
			toRender += style.Text(word[len(input)+1:])
		}
	}

	// Update the word to render.
	p.wordsToRender[p.currentWord] = toRender

	return strings.Join(p.wordsToRender, " ")
}

func (p PlayHandler) isFinished() bool {
	return p.currentWord == len(p.words)
}

func (p PlayHandler) isEmptyInput() bool {
	return len(p.textInput.Value()) == 0
}

// Styles the provided text depending on the user input.
// Returns the text styled.
func (p PlayHandler) textStatus(text string) string {
	current := p.words[p.currentWord]
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

func (p *PlayHandler) Render() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		p.renderWords(),
		p.textInput.View(),
	) + "\n"
}
