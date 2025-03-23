package utils

import (
	"strings"
)

func NormalizeInnerSpaces(s string) string {
	// Divide en palabras, quita los vacíos, y las une con un solo espacio
	parts := strings.Fields(s) // hace split y elimina espacios duplicados automáticamente
	return strings.Join(parts, " ")
}

func NormalizeText(s string) string {
	return NormalizeInnerSpaces(strings.ToLower(strings.TrimSpace(s)))
}
