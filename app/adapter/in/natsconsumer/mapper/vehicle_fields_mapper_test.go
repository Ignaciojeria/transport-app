package mapper

import (
	"testing"
)

func TestVehicleFieldMapper_Map(t *testing.T) {
	mapper := NewVehicleFieldMapper()

	tests := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]string
	}{
		{
			name: "Mapeo completo con claves canónicas",
			input: map[string]interface{}{
				"id":                     "vehicle1",
				"weight":                 "1000",
				"volume":                 "500",
				"insurance":              "100",
				"startLocationLatitude":  "40.4168",
				"startLocationLongitude": "-3.7038",
				"endLocationLatitude":    "41.3851",
				"endLocationLongitude":   "2.1734",
				"maxPackageQuantity":     "10",
			},
			expected: map[string]string{
				"endLocationLatitude":    "endLocationLatitude",
				"endLocationLongitude":   "endLocationLongitude",
				"id":                     "id",
				"insurance":              "insurance",
				"startLocationLatitude":  "startLocationLatitude",
				"startLocationLongitude": "startLocationLongitude",
				"volume":                 "volume",
				"weight":                 "weight",
				"maxPackageQuantity":     "maxPackageQuantity",
			},
		},
		{
			name: "Mapeo con claves alternativas",
			input: map[string]interface{}{
				"vehicle_id":   "vehicle1",
				"peso":         "1000",
				"volumen":      "500",
				"seguro":       "100",
				"origen_lat":   "40.4168",
				"origen_lon":   "-3.7038",
				"destino_lat":  "41.3851",
				"destino_lon":  "2.1734",
				"max_paquetes": "10",
			},
			expected: map[string]string{
				"endLocationLatitude":    "endLocationLatitude",
				"endLocationLongitude":   "endLocationLongitude",
				"id":                     "id",
				"insurance":              "insurance",
				"startLocationLatitude":  "startLocationLatitude",
				"startLocationLongitude": "startLocationLongitude",
				"volume":                 "volume",
				"weight":                 "weight",
				"maxPackageQuantity":     "maxPackageQuantity",
			},
		},
		{
			name: "Mapeo con claves mixtas",
			input: map[string]interface{}{
				"id":                     "vehicle1",
				"weight":                 "1000",
				"volume":                 "500",
				"insurance":              "100",
				"startLocationLatitude":  "40.4168",
				"startLocationLongitude": "-3.7038",
				"endLocationLatitude":    "41.3851",
				"endLocationLongitude":   "2.1734",
				"max_paquetes":           "10",
			},
			expected: map[string]string{
				"endLocationLatitude":    "endLocationLatitude",
				"endLocationLongitude":   "endLocationLongitude",
				"id":                     "id",
				"insurance":              "insurance",
				"startLocationLatitude":  "startLocationLatitude",
				"startLocationLongitude": "startLocationLongitude",
				"volume":                 "volume",
				"weight":                 "weight",
				"maxPackageQuantity":     "maxPackageQuantity",
			},
		},
		{
			name: "Mapeo con claves que no existen",
			input: map[string]interface{}{
				"id":                     "vehicle1",
				"weight":                 "1000",
				"volume":                 "500",
				"insurance":              "100",
				"startLocationLatitude":  "40.4168",
				"startLocationLongitude": "-3.7038",
				"endLocationLatitude":    "41.3851",
				"endLocationLongitude":   "2.1734",
				"unknown_field":          "value",
			},
			expected: map[string]string{
				"endLocationLatitude":    "endLocationLatitude",
				"endLocationLongitude":   "endLocationLongitude",
				"id":                     "id",
				"insurance":              "insurance",
				"startLocationLatitude":  "startLocationLatitude",
				"startLocationLongitude": "startLocationLongitude",
				"volume":                 "volume",
				"weight":                 "weight",
				"maxPackageQuantity":     "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.Map(tt.input)

			// Verificar que todas las claves esperadas estén presentes
			for expectedKey := range tt.expected {
				if _, exists := result[expectedKey]; !exists {
					t.Errorf("Clave esperada '%s' no encontrada en el resultado", expectedKey)
				}
			}

			// Verificar que no haya claves extra en el resultado
			for resultKey := range result {
				if _, exists := tt.expected[resultKey]; !exists {
					t.Errorf("Clave inesperada '%s' encontrada en el resultado", resultKey)
				}
			}

			// Verificar que las claves canónicas estén presentes (pueden tener valores vacíos si no se encontró la clave original)
			for expectedKey := range tt.expected {
				if _, exists := result[expectedKey]; !exists {
					t.Errorf("Clave canónica '%s' no encontrada en el resultado", expectedKey)
				}
			}
		})
	}
}

func TestVehicleFieldMapper_Map_MaxPackageQuantity(t *testing.T) {
	mapper := NewVehicleFieldMapper()

	// Test específico para maxPackageQuantity con diferentes variaciones
	testCases := []string{
		"maxPackageQuantity",
		"max_package_quantity",
		"max_paquetes",
		"max_paquete",
		"max_paquetes_maximo",
		"max_paquete_maximo",
	}

	for _, testKey := range testCases {
		t.Run("Test_"+testKey, func(t *testing.T) {
			input := map[string]interface{}{
				testKey: "15",
			}

			result := mapper.Map(input)

			if result["maxPackageQuantity"] != testKey {
				t.Errorf("Para la clave de entrada '%s', se esperaba que se mapeara a 'maxPackageQuantity', pero se obtuvo '%s'",
					testKey, result["maxPackageQuantity"])
			}
		})
	}
}
