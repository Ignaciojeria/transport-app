package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// NormalizeText normaliza el texto eliminando acentos, signos de puntuación,
// reemplazando "ñ" por "n" y convirtiéndolo en minúsculas sin modificar los espacios.
func NormalizeText(text string) string {
	// Eliminar acentos y caracteres especiales
	t := norm.NFD.String(text) // Descompone caracteres con acentos (ej: "á" → "á")
	var sb strings.Builder
	for _, r := range t {
		if unicode.IsMark(r) { // IsMark detecta acentos y los elimina
			continue
		}
		sb.WriteRune(r)
	}
	text = sb.String()

	// Convertir a minúsculas
	text = strings.ToLower(text)

	// Reemplazar "ñ" por "n"
	text = strings.ReplaceAll(text, "ñ", "n")
	text = strings.ReplaceAll(text, "Ñ", "N")

	// Eliminar signos de puntuación pero mantener espacios
	re := regexp.MustCompile(`[^\w\s]`) // Permite solo letras, números y espacios
	text = re.ReplaceAllString(text, "")

	// Eliminar espacios en blanco al inicio y final
	return strings.TrimSpace(text)
}
