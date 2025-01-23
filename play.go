package main

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vieitesss/gotype/keymaps"
	"github.com/vieitesss/gotype/style"
)

// Used when playing, typing.
type PlayHandler struct {
	viewport      PlayViewport
	keys          keymaps.PlayKeyMaps
	help          help.Model
	textInput     textinput.Model
	timer         timer.Model
	words         []string
	wordsToRender []string
	lastIncorrect []string
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
		keys:      keymaps.PlayKeys(),
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
		tea.WindowSize(),
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
		startWord := p.viewport.FirstWordFromLine(p.viewport.currentLine)
		p.viewport.SetContent(strings.Join(p.wordsToRender[startWord:], " "))

		return p, nil

	case NextWordMsg:
		if p.currentWord+1 == len(p.words) {
			return p, updateStatus(Quit)
		}

		p.currentWord++
		p.textInput.Reset()

		return p, nil

	case PrevWordMsg:
		p.currentWord--
		lastLen := len(p.lastIncorrect)
		p.textInput.SetValue(p.lastIncorrect[lastLen-1])
		p.lastIncorrect = p.lastIncorrect[:lastLen-1]

		p.viewport.ToPrevWord()

	case TextToWriteMsg:
		p.words = msg
		p.wordsToRender = style.InitialWordsStyling(p.words)
		p.viewport.SetWordsPerLine(p.words)

		return p, nil

	case tea.WindowSizeMsg:
		if p.viewport.ready {
			p.viewport.UpdateFrame(msg.Width, p.words)
		} else {
			p.viewport = NewPlayViewport(msg.Width, strings.Join(p.wordsToRender, " "))
		}

		p.textInput.PromptStyle = style.NormalStyle.PaddingLeft(p.viewport.margin)

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
			var cmds []tea.Cmd
			msg := p.updateCurrentWordStyle(p.currentWord, false)()

			switch updated := msg.(type) {
			case UpdatedWordToRenderMsg:
				p.wordsToRender[p.currentWord] = string(updated)

			default:
				panic("[ERROR] play.go:Update updated type should be UpdatedWordToRenderMsg")
			}

			p.checkIfIncorrect(p.currentWord)

			// Go to the next word.
			cmds = append(cmds, p.nextWord)

			// TODO: show the results.

			// No more words
			if p.currentWord == len(p.words) {
				cmds = append(cmds, updateStatus(Quit))
			}

			p.viewport.ToNextWord()

			return p, tea.Batch(cmds...)

		case tea.KeyBackspace:
			if len(p.textInput.Value()) == 0 && len(p.lastIncorrect) > 0 {
				p.wordsToRender[p.currentWord] = style.Text(p.words[p.currentWord])

				return p, p.prevWord
			}

		default:
			if !p.started {
				p.started = true
				cmds = append(cmds, p.timer.Init())
			}
		}
	}

	if len(p.words) > 0 {
		cmds = append(cmds, p.updateCurrentWordStyle(p.currentWord, true))
	}

	var cmd tea.Cmd
	p.textInput, cmd = p.textInput.Update(msg)
	cmds = append(cmds, cmd)
	p.viewport.model, cmd = p.viewport.model.Update(msg)
	cmds = append(cmds, cmd)

	return p, tea.Batch(cmds...)
}

func (p *PlayHandler) checkIfIncorrect(index int) {
	current := p.textInput.Value()
	if current != p.words[index] {
		// Keep track of last incorrect words.
		p.lastIncorrect = append(p.lastIncorrect, current)
	} else if len(p.lastIncorrect) > 0 {
		// Empty lastIncorrect if a correct word is written.
		p.lastIncorrect = nil
	}
}

func (p *PlayHandler) updateCurrentWordStyle(index int, addCursor bool) tea.Cmd {
	return func() tea.Msg {
		styled := style.CompareWithStyle(
			p.textInput.Value(),
			p.words[index],
			addCursor,
		)

		return UpdatedWordToRenderMsg(styled)
	}
}

func (p *PlayHandler) nextWord() tea.Msg {
	return NextWordMsg{}
}

func (p *PlayHandler) prevWord() tea.Msg {
	return PrevWordMsg{}
}

func (p PlayHandler) Render() string {
	s := ""

	// Timer
	s += strings.Repeat(" ", p.viewport.margin)
	if !p.timer.Timedout() {
		s += p.timer.View() + "\n"
	} else {
		s += "Time is over!\n"
	}

	// Words
	s += p.viewport.model.View() + "\n"

	// Input
	s += "\n" + p.textInput.View() + "\n"

	// Help
	s += strings.Repeat("\n", 3) + "  " + p.help.View(p.keys)

	return s
}
