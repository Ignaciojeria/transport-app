package schema

import "google.golang.org/genai"

// GetSideSchema define la estructura para un acompañamiento (ej. extra de queso).
func GetSideSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
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
		},
		Required: []string{"name", "pricing"},
	}
}
