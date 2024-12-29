package style

import (
	"github.com/charmbracelet/lipgloss"
)

const white = lipgloss.Color("#ffffff")
const green = lipgloss.Color("#06d6a0")
const red = lipgloss.Color("#ef476f")

var normalStyle = lipgloss.NewStyle()
var boldStyle = lipgloss.NewStyle().Bold(true)

var textStyle = boldStyle.Foreground(white)
var correctStyle = boldStyle.Foreground(green)
var incorrectStyle = boldStyle.Foreground(red)
var selectorStyle = normalStyle.Foreground(green)

func Text(text string) string {
	return textStyle.Render(text)
}

func Selector(text string) string {
	return correctStyle.Render(text)
}

func Correct(text string) string {
	return correctStyle.Render(text)
}

func Incorrect(text string) string {
	return incorrectStyle.Render(text)
}
