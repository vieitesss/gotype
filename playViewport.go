package main

import (
	"math"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/vieitesss/gotype/style"
)

type PlayViewport struct {
	model        viewport.Model
	ready        bool
	margin       int
	lineWidth    int
	wordsPerLine []int
	wordsWritten int
	currentLine  int
}

var viHeigth = 3

func NewPlayViewport(width int, wordsToRender string) PlayViewport {
	p := PlayViewport{}

	margin := p.CalculateMargin(width)

	vi := viewport.New(width, viHeigth)
	vi.Style = style.NormalStyle.Padding(0, margin)
	vi.SetContent(wordsToRender)

	p.model = vi
	p.margin = margin
	p.lineWidth = width - margin*2
	p.ready = true

	return p
}

func (pv *PlayViewport) UpdateFrame(width int, words []string) {
	pv.model.Width = width
	pv.margin = pv.CalculateMargin(width)
	pv.model.Style = style.NormalStyle.Padding(0, pv.margin)
	pv.SetWordsPerLine(words)
	pv.lineWidth = width - pv.margin*2
}

func (PlayViewport) CalculateMargin(width int) int {
	return int(math.Floor(float64(width) * 0.3 / 2))
}

func (pv *PlayViewport) SetContent(content string) {
	if pv.ready {
		pv.model.SetContent(content)
	}
}

func (pv *PlayViewport) ToPrevWord() {
	if pv.wordsWritten > 0 {
		pv.wordsWritten--
	} else {
		pv.currentLine--
		pv.wordsWritten = pv.wordsPerLine[pv.currentLine] - 1
	}
}

func (pv *PlayViewport) ToNextWord() {
	pv.wordsWritten += 1
	if pv.wordsWritten == pv.wordsPerLine[pv.currentLine] {
		pv.currentLine += 1
		pv.wordsWritten = 0
	}
}

func (pv *PlayViewport) SetWordsPerLine(words []string) {
	if pv.lineWidth < 1 {
		return
	}

	chars := 0
	lineWords := 0
	pv.wordsPerLine = make([]int, 0)

	for _, w := range words {
		wordLen := len([]rune(w))
		chars += wordLen

		if chars > pv.lineWidth {
			pv.wordsPerLine = append(pv.wordsPerLine, lineWords)
			lineWords = 1
			chars = wordLen
		} else {
			lineWords += 1
		}
		chars += 1 // The space.
	}

	pv.wordsPerLine = append(pv.wordsPerLine, lineWords)
}

func (pv PlayViewport) FirstWordFromLine(line int) int {
	if line == 0 && line == 1 {
		return 0
	}

	index := 0

	for i := 0; i < line-1; i++ {
		index += pv.wordsPerLine[i]
	}

	return index
}
