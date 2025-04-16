package utils

import (
	"testing"
)

func TestNormalizeInnerSpaces(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"   hola    mundo   ", "hola mundo"},
		{"uno dos   tres", "uno dos tres"},
		{"   espacio   al   inicio ", "espacio al inicio"},
		{"sin    espacios     extras", "sin espacios extras"},
		{"", ""},
		{"   ", ""},
	}

	for _, test := range tests {
		result := NormalizeInnerSpaces(test.input)
		if result != test.expected {
			t.Errorf("NormalizeInnerSpaces(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}

func TestNormalizeText(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Espacios, minúsculas
		{"  AVENIDA   PROVIDENCIA   1234 ", "avenida providencia 1234"},
		{"Calle    los  Castaños", "calle los castaños"},
		{"   Ñuñoa   ", "ñuñoa"},
		{"", ""},
		{"   ", ""},

		// Acentos
		{"PéREz     de  Valdivia", "perez de valdivia"},
		{"Región Metropolitana", "region metropolitana"},
		{"José María", "jose maria"},
		{"Camión", "camion"},

		// Puntuación
		{"¡Hola, mundo!", "hola mundo"},
		{"¿Dónde estás?", "donde estas"},
		{"Av. Las Condes Nº1234, Dpto. 5B", "av las condes nº1234 dpto 5b"}, // conserva nº si no querés filtrar letras especiales

		// Símbolos no alfabéticos (pueden decidirse según requerimiento)
		{"Correo: usuario@example.com", "correo usuarioexamplecom"},
		{"$100.000 CLP", "100000 clp"},
	}

	for _, test := range tests {
		result := NormalizeText(test.input)
		if result != test.expected {
			t.Errorf("NormalizeText(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
