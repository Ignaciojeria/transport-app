package schema

import "google.golang.org/genai"

// GetMenuItemSchema define la estructura para un producto individual del menú.
func GetMenuItemSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"id": {
				Type:        genai.TypeString,
				Description: "Identificador semántico único del item en formato kebab-case (ej. 'empanadas-pino', 'pizza-margherita').",
			},
			"title": {
				Type:        genai.TypeObject,
				Description: "Título del producto en formato multiidioma con base y traducciones.",
				Properties: map[string]*genai.Schema{
					"base": {
						Type:        genai.TypeString,
						Description: "Texto base del título (idioma principal, generalmente español).",
					},
					"languages": {
						Type:        genai.TypeObject,
						Description: "Traducciones del título en diferentes idiomas.",
						Properties: map[string]*genai.Schema{
							"es": {Type: genai.TypeString, Description: "Título en español."},
							"en": {Type: genai.TypeString, Description: "Título en inglés."},
							"pt": {Type: genai.TypeString, Description: "Título en portugués."},
						},
					},
				},
				Required: []string{"base", "languages"},
			},
			"description": {
				Type:        genai.TypeObject,
				Description: "Descripción del producto en formato multiidioma con base y traducciones. Opcional.",
				Properties: map[string]*genai.Schema{
					"base": {
						Type:        genai.TypeString,
						Description: "Texto base de la descripción (idioma principal, generalmente español).",
					},
					"languages": {
						Type:        genai.TypeObject,
						Description: "Traducciones de la descripción en diferentes idiomas.",
						Properties: map[string]*genai.Schema{
							"es": {Type: genai.TypeString, Description: "Descripción en español."},
							"en": {Type: genai.TypeString, Description: "Descripción en inglés."},
							"pt": {Type: genai.TypeString, Description: "Descripción en portugués."},
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
			"sides": {
				Type: genai.TypeArray,
				// ¡Aquí se utiliza el Schema del Side!
				Items:       GetSideSchema(),
				Description: "Lista opcional de acompañamientos o variaciones con pricing (ej. tamaño grande, extra de salsa).",
			},
			"pricing": {
				Type:        genai.TypeObject,
				Description: "Estructura de pricing del ítem del menú.",
				Properties:  GetPricingSchema().Properties,
				Required:    GetPricingSchema().Required,
			},
			"photoUrl": {
				Type:        genai.TypeString,
				Description: "URL opcional de la imagen del producto. Debe ser una URL pública accesible.",
			},
			"station": {
				Type:        genai.TypeString,
				Enum:        []string{"KITCHEN", "BAR"},
				Description: "Estación que prepara este ítem: KITCHEN (cocina) o BAR (bar). Opcional.",
			},
		},
		Required: []string{"id", "title", "pricing"},
	}
}
