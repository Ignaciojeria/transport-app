package schema

import "google.golang.org/genai"

// GetSideSchema define la estructura para un acompañamiento (ej. extra de queso).
func GetSideSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"id": {
				Type:        genai.TypeString,
				Description: "Identificador semántico único del side en formato kebab-case (ej. 'papas-fritas', 'tamaño-grande').",
			},
			"name": {
				Type:        genai.TypeString,
				Description: "Nombre del acompañamiento (ej. 'Extra de tocino').",
			},
			"pricing": {
				Type:        genai.TypeObject,
				Description: "Estructura de pricing del acompañamiento.",
				Properties:  GetPricingSchema().Properties,
				Required:    GetPricingSchema().Required,
			},
			"photoUrl": {
				Type:        genai.TypeString,
				Description: "URL opcional de la imagen del acompañamiento. Debe ser una URL pública accesible.",
			},
			"station": {
				Type:        genai.TypeString,
				Enum:        []string{"KITCHEN", "BAR"},
				Description: "Estación que prepara este acompañamiento: KITCHEN (cocina) o BAR (bar). Opcional.",
			},
		},
		Required: []string{"id", "name", "pricing"},
	}
}
