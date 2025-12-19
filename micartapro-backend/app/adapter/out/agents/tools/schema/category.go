package schema

import "google.golang.org/genai"

// GetMenuCategorySchema define la estructura para una categoría del menú (ej. Platos Fuertes).
func GetMenuCategorySchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"title": {Type: genai.TypeString, Description: "Título de la categoría (ej. 'Hamburguesas Especiales')."},
			"items": {
				Type: genai.TypeArray,
				// ¡Aquí se utiliza el Schema del Ítem!
				Items:       GetMenuItemSchema(),
				Description: "Lista de todos los productos (ítems) dentro de esta categoría.",
			},
		},
		Required: []string{"title", "items"},
	}
}
