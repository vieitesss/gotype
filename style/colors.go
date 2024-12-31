package style

import (
	"github.com/charmbracelet/lipgloss"
)

const white = lipgloss.Color("#dddddd")
const green = lipgloss.Color("#accca4")
const red = lipgloss.Color("#bf9ac1")

var normalStyle = lipgloss.NewStyle()
var boldStyle = lipgloss.NewStyle().Bold(true)

var textStyle = boldStyle.Foreground(white)
var correctStyle = boldStyle.Foreground(green)
var incorrectStyle = boldStyle.Foreground(red).Underline(true)
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
