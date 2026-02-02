package schema

import "google.golang.org/genai"

// GetMenuCategorySchema define la estructura para una categoría del menú (ej. Platos Fuertes).
func GetMenuCategorySchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"title": {
				Type:        genai.TypeObject,
				Description: "Título de la categoría en formato multiidioma con base y traducciones.",
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
