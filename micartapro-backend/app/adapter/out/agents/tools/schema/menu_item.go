package schema

import "google.golang.org/genai"

// GetMenuItemSchema define la estructura para un producto individual del menú.
func GetMenuItemSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"title":       {Type: genai.TypeString, Description: "Nombre del producto (ej. 'Hamburguesa Clásica')."},
			"description": {Type: genai.TypeString, Description: "Breve descripción del producto. Puede estar vacío."},
			"sides": {
				Type: genai.TypeArray,
				// ¡Aquí se utiliza el Schema del Side!
				Items:       GetSideSchema(),
				Description: "Lista opcional de acompañamientos o variaciones con precio (ej. tamaño grande, extra de salsa).",
			},
			"price": {Type: genai.TypeNumber, Description: "Precio principal del ítem si no tiene sides."},
		},
		Required: []string{"title"},
	}
}
