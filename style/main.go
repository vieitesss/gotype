package style

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func MainMenuOptionStyling(option string) string {
	parts := strings.Split(option, " ")

	if len(parts) > 1 {
		return Selector(parts[0]) + " " + Text(parts[1])
	}

	return "  " + Text(parts[0])
}

func InitialWordsStyling(words []string) []string {
	styled := make([]string, len(words))

	for i, w := range words {
		styled[i] = Text(w)
	}

	return styled
}

// Source is what you guess.
// Target is what you want.
func CompareWithStyle(source, target string, addCursor bool) string {
	toRender := ""
	sourceRune, targetRune := []rune(source), []rune(target)
	sourceLen, targetLen := len(sourceRune), len(targetRune)
	minLen := sourceLen

	if sourceLen > targetLen {
		minLen = targetLen
	}

	for i := 0; i < minLen; i++ {
		if sourceRune[i] == targetRune[i] {
			toRender += styleWithFunc(targetRune, i, i+1, Correct)
		} else {
			toRender += styleWithFunc(targetRune, i, i+1, Incorrect)
		}
	}

	if sourceLen == targetLen {
		return toRender
	}

	if sourceLen > targetLen {
		return toRender + styleWithFunc(sourceRune, targetLen, sourceLen, Incorrect)
	}

	if !addCursor {
		return toRender + styleWithFunc(targetRune, sourceLen, targetLen, Text)
	}

	toRender += styleWithFunc(targetRune, sourceLen, sourceLen+1, Cursor)

	if targetLen > sourceLen+1 {
		toRender += styleWithFunc(targetRune, sourceLen+1, targetLen, Text)
	}

	return toRender
}

// Returns the text styled and the number of characters checked.
func styleWithFunc(text []rune, start, finish int, callback func(string) string) string {
	if start >= finish {
		panic(fmt.Sprintf("[ERROR] main.go:styleWithFunc - %d (start) should be",
			"smaller than %d (finish).", start, finish))
	}

	toRender := text[start:finish]

	return callback(string(toRender))
}
