package style

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	white = lipgloss.Color("#dddddd")
	green = lipgloss.Color("#accca4")
	red   = lipgloss.Color("#bf9ac1")
	lightGrey = lipgloss.Color("#333333")
)

var (
	NormalStyle = lipgloss.NewStyle()
	BoldStyle   = lipgloss.NewStyle().Bold(true)

	TextStyle      = BoldStyle.Foreground(white)
	MenuStyle      = TextStyle.Background(lightGrey).PaddingLeft(1).PaddingRight(1)
	CorrectStyle   = BoldStyle.Foreground(green)
	IncorrectStyle = BoldStyle.Foreground(red).Underline(true)
	SelectorStyle  = NormalStyle.Foreground(green)
)

func Text(text string) string {
	return TextStyle.Render(text)
}

func Selector(text string) string {
	return CorrectStyle.Render(text)
}

func Correct(text string) string {
	return CorrectStyle.Render(text)
}

func Incorrect(text string) string {
	return IncorrectStyle.Render(text)
}

func Cursor(text string) string {
	return TextStyle.Underline(true).Render(text)
}
