package utils

import (
	"strings"
)

var replacements = map[rune]rune{
	'á': 'a', 'é': 'e', 'í': 'i', 'ó': 'o', 'ú': 'u',
	'Á': 'a', 'É': 'e', 'Í': 'i', 'Ó': 'o', 'Ú': 'u',
	'ä': 'a', 'ë': 'e', 'ï': 'i', 'ö': 'o', 'ü': 'u',
	'Ä': 'a', 'Ë': 'e', 'Ï': 'i', 'Ö': 'o', 'Ü': 'u',
	'à': 'a', 'è': 'e', 'ì': 'i', 'ò': 'o', 'ù': 'u',
	'À': 'a', 'È': 'e', 'Ì': 'i', 'Ò': 'o', 'Ù': 'u',
}

func NormalizeInnerSpaces(s string) string {
	parts := strings.Fields(s)
	return strings.Join(parts, " ")
}

func NormalizeText(s string) string {
	s = NormalizeInnerSpaces(strings.ToLower(strings.TrimSpace(s)))

	var b strings.Builder
	for _, r := range s {
		if repl, ok := replacements[r]; ok {
			b.WriteRune(repl)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
