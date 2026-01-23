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
				Type:        genai.TypeString,
				Description: "Nombre del producto (ej. 'Hamburguesa Clásica').",
			},
			"description": {
				Type:        genai.TypeString,
				Description: "Breve descripción del producto. Puede estar vacío.",
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
		},
		Required: []string{"id", "title", "pricing"},
	}
}
