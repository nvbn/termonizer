package utils

func FindLineStart(text string, pos int) int {
	for pos := min(pos, len(text)-1); pos > 0; pos-- {
		if len(text) > pos && text[pos-1] == '\n' {
			return pos
		}
	}

	return 0
}

func FindLineEnd(text string, pos int) int {
	if len(text) == 0 {
		return 0
	}

	if pos >= len(text)-1 {
		return len(text) - 1
	}

	for ; pos < len(text); pos++ {
		if text[pos] == '\n' {
			return pos
		}
	}

	return len(text) - 1
}
