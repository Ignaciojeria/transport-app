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
		{"  AVENIDA   PROVIDENCIA   1234 ", "avenida providencia 1234"},
		{"Calle    los  Castaños", "calle los castaños"},
		{"   Ñuñoa   ", "ñuñoa"},
		{"", ""},
		{"   ", ""},
		{"PéREz     de  Valdivia", "pérez de valdivia"},
	}

	for _, test := range tests {
		result := NormalizeText(test.input)
		if result != test.expected {
			t.Errorf("NormalizeText(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
