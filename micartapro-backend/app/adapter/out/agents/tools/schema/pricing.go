package schema

import "google.golang.org/genai"

// GetPricingSchema define la estructura para el pricing de un ítem o side.
func GetPricingSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"mode": {
				Type:        genai.TypeString,
				Description: "Modo de pricing: 'UNIT' (por unidad), 'WEIGHT' (por peso), 'VOLUME' (por volumen), 'LENGTH' (por longitud), 'AREA' (por área).",
			},
			"unit": {
				Type:        genai.TypeString,
				Description: "Unidad de medida: 'EACH' (unidad), 'GRAM' (gramo), 'KILOGRAM' (kg), 'MILLILITER' (ml), 'LITER' (l), 'METER' (m), 'SQUARE_METER' (m²).",
			},
			"pricePerUnit": {
				Type:        genai.TypeNumber,
				Description: "Precio por unidad de medida. Usa la moneda del negocio (businessInfo.currency): si es CLP usa número entero (sin decimales); si es USD o BRL se permiten decimales (ej: 9.99).",
			},
			"baseUnit": {
				Type:        genai.TypeNumber,
				Description: "Unidad base opcional (ej: 100 => precio por 100 gramos). Puede omitirse si no aplica.",
			},
		},
		Required: []string{"mode", "unit", "pricePerUnit"},
	}
}
