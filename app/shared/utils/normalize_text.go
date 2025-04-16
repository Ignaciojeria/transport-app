package utils

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func NormalizeInnerSpaces(s string) string {
	parts := strings.Fields(s)
	return strings.Join(parts, " ")
}

func NormalizeText(s string) string {
	// 1. Minusculizar y quitar espacios extra
	s = NormalizeInnerSpaces(strings.ToLower(strings.TrimSpace(s)))

	// 2. Normalizar Unicode (NFD) para separar tildes
	nfd := norm.NFD.String(s)

	// 3. Eliminar tildes y puntuación, pero dejar la ñ
	var b strings.Builder
	for _, r := range nfd {
		if unicode.Is(unicode.Mn, r) {
			continue // Elimina marcas de acento
		}
		if unicode.IsPunct(r) {
			continue // Elimina puntuación (.,;:!? etc.)
		}
		b.WriteRune(r)
	}
	return b.String()
}
