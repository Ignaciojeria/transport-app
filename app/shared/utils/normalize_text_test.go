package utils

import (
	"testing"
)

func TestNormalizeText(t *testing.T) {
	cases := map[string][]struct {
		input    string
		expected string
	}{
		"espacios y minúsculas": {
			{"  AVENIDA   PROVIDENCIA   1234 ", "avenida providencia 1234"},
			{"Calle    los  Castaños", "calle los castaños"},
			{"   Ñuñoa   ", "ñuñoa"},
			{"", ""},
			{"   ", ""},
		},
		"acentos": {
			{"PéREz     de  Valdivia", "perez de valdivia"},
			{"Región Metropolitana", "region metropolitana"},
			{"José María", "jose maria"},
			{"Camión", "camion"},
		},
		"puntuación": {
			{"¡Hola, mundo!", "¡hola, mundo!"},
			{"¿Dónde estás?", "¿donde estas?"},
			{"Av. Las Condes Nº1234, Dpto. 5B", "av. las condes nº1234, dpto. 5b"},
		},
		"símbolos y otros": {
			{"Correo: usuario@example.com", "correo: usuario@example.com"},
			{"$100.000 CLP", "$100.000 clp"},
		},
	}

	for category, tests := range cases {
		t.Run(category, func(t *testing.T) {
			for _, test := range tests {
				result := NormalizeText(test.input)
				if result != test.expected {
					t.Errorf("NormalizeText(%q) = %q; want %q", test.input, result, test.expected)
				}
			}
		})
	}
}
