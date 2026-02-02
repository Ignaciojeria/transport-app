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
				Type:        genai.TypeObject,
				Description: "Nombre del acompañamiento en formato multiidioma con base y traducciones.",
				Properties: map[string]*genai.Schema{
					"base": {
						Type:        genai.TypeString,
						Description: "Texto base del nombre (idioma principal, generalmente español).",
					},
					"languages": {
						Type:        genai.TypeObject,
						Description: "Traducciones del nombre en diferentes idiomas.",
						Properties: map[string]*genai.Schema{
							"es": {Type: genai.TypeString, Description: "Nombre en español."},
							"en": {Type: genai.TypeString, Description: "Nombre en inglés."},
							"pt": {Type: genai.TypeString, Description: "Nombre en portugués."},
						},
					},
				},
				Required: []string{"base", "languages"},
			},
			"foodAttributes": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeString,
					Enum: []string{"GLUTEN", "SEAFOOD", "NUTS", "DAIRY", "EGGS", "SOY", "VEGAN", "VEGETARIAN", "SPICY", "ALCOHOL"},
				},
				Description: "Atributos alimentarios opcionales: alérgenos (GLUTEN, SEAFOOD, NUTS, DAIRY, EGGS, SOY), dieta (VEGAN, VEGETARIAN), SPICY, ALCOHOL.",
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
