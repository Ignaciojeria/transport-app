package mapper

import (
	"strings"
	"unicode"
)

func deaccent(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		case 'á', 'à', 'ä', 'â':
			r = 'a'
		case 'é', 'è', 'ë', 'ê':
			r = 'e'
		case 'í', 'ì', 'ï', 'î':
			r = 'i'
		case 'ó', 'ò', 'ö', 'ô':
			r = 'o'
		case 'ú', 'ù', 'ü', 'û':
			r = 'u'
			// case 'ñ':
			//     r = 'n'  // Preservar ñ como letra válida
		}
		if r == 0 || !unicode.IsControl(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}
