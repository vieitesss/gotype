package main

import (
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/bubbles/viewport"
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
	viewport      viewport.Model
	timer         timer.Model
	words         []string
	wordsToRender []string
	lastIncorrect []string
	currentWord   int
	started       bool
	vpReady       bool
	seconds       int
	margin        int
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
		if p.vpReady {
			p.viewport.SetContent(strings.Join(p.wordsToRender, " "))
		}

		return p, nil

	case NextWordMsg:
		p.currentWord++
		p.textInput.Reset()

		return p, nil

	case PrevWordMsg:
		p.currentWord--
		lastLen := len(p.lastIncorrect)
		p.textInput.SetValue(p.lastIncorrect[lastLen-1])
		p.lastIncorrect = p.lastIncorrect[:lastLen-1]

	case TextToWriteMsg:
		p.words = msg
		p.wordsToRender = style.InitialWordsStyling(p.words)

		return p, nil

	case tea.WindowSizeMsg:
		p.margin = int(math.Floor(float64(msg.Width) * 0.3 / 2))
		p.textInput.PromptStyle = style.NormalStyle.PaddingLeft(p.margin)

		if !p.vpReady {
			p.viewport = viewport.New(msg.Width, 4)
			p.viewport.Style = style.NormalStyle.Padding(0, p.margin)
			p.viewport.SetContent(strings.Join(p.wordsToRender, " "))
			p.viewport.YPosition = 10
			p.vpReady = true
		} else {
			p.viewport.Width = msg.Width
		}

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
			msg := p.updateCurrentWord(p.currentWord, false)()
			switch updated := msg.(type) {
			case UpdatedWordToRenderMsg:
				p.wordsToRender[p.currentWord] = string(updated)

			default:
				panic("[ERROR] play.go:Update updated type should be UpdatedWordToRenderMsg")
			}

			current := p.textInput.Value()
			if current != p.words[p.currentWord] {
				// Keep track of last incorrect words.
				p.lastIncorrect = append(p.lastIncorrect, current)
			} else if len(p.lastIncorrect) > 0 {
				// Empty lastIncorrect if a correct word is written.
				p.lastIncorrect = nil
			}

			// Go to the next word.
			cmds = append(cmds, p.nextWord)

			// TODO: show the results.

			// No more words
			if p.currentWord == len(p.words) {
				cmds = append(cmds, updateStatus(Quit))
			}

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
		cmds = append(cmds, p.updateCurrentWord(p.currentWord, true))
	}

	var cmd tea.Cmd
	p.textInput, cmd = p.textInput.Update(msg)
	cmds = append(cmds, cmd)
	p.viewport, cmd = p.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return p, tea.Batch(cmds...)
}

func (p *PlayHandler) updateCurrentWord(index int, addCursor bool) tea.Cmd {
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
	if !p.timer.Timedout() {
		s += strings.Repeat(" ", p.margin) + p.timer.View() + "\n"
	} else {
		s += "Time is over!\n"
	}

	// Words
	s += p.viewport.View() + "\n"

	// Input
	s += "\n" + p.textInput.View() + "\n"

	// Help
	s += strings.Repeat("\n", 3) + "  " + p.help.View(p.keys)

	return s
}
