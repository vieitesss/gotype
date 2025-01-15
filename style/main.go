package style

func InitialWordsStyling(words []string) []string {
	styled := make([]string, len(words))

	for i, w := range words {
		styled[i] = Text(w)
	}

	return styled
}

func CompareWithStyle(source, target string, addCursor bool) string {
	toRender := ""
	sourceLen := len(source)
	targetLen := len(target)
	minLen := sourceLen
	if sourceLen > targetLen {
		minLen = targetLen
	}

	for i := 0; i < minLen; i++ {
		if source[i] == target[i] {
			toRender += Correct(string(target[i]))
		} else {
			toRender += Incorrect(string(target[i]))
		}
	}

	if sourceLen == targetLen {
		return toRender
	}

	if sourceLen > targetLen {
		toRender += Incorrect(source[targetLen:])
		return toRender
	}

	if addCursor {
		toRender += CurrentChar(string(target[sourceLen]))
		if targetLen > sourceLen+1 {
			toRender += Text(target[sourceLen+1:])
		}
	} else {
		toRender += Text(target[sourceLen:])
	}

	return toRender
}
