package go_printer

import "strings"

func splitLines(text string, size int) []string {
	lines := make([]string, 0)

	if len(text) == 0 {
		return []string{""}
	}

	for len(text) > size {
		pos := size
		for text[pos] != ' ' && pos > 0 {
			pos--
		}
		if pos == 0 {
			pos = size
		}
		lines = append(lines, strings.TrimRight(text[0:pos], " "))
		text = strings.TrimLeft(text[pos:], " ")
	}

	if len(text) > 0 {
		lines = append(lines, text)
	}

	return lines
}

func splitTextIntoLines(text string, size int) []string {
	lines := make([]string, 0)

	for _, line := range strings.Split(text, "\n") {
		lines = append(lines, splitLines(line, size)...)
	}

	return lines
}
