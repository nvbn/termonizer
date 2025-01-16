package utils

func FindLineStart(text string, pos int) int {
	if text == "" {
		return 0
	}

	pos = min(pos, len(text)-1)
	if text[pos] == '\n' {
		return pos + 1 // special case, could confuse on empty line
	}

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
