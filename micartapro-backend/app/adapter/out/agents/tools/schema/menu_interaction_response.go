package schema

import "google.golang.org/genai"

// MenuInteractionResponse define la estructura raíz para la respuesta final del menú.
func MenuInteractionResponse() *genai.Schema {
	contactSchema := GetContactSchema()
	menuCategorySchema := GetMenuCategorySchema()

	// Nota: Reestructuramos un poco el contacto/horario para ser más estándar
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"business": {
				Type:        genai.TypeObject,
				Properties:  contactSchema.Properties, // Reutilizamos las propiedades del ContactSchema
				Description: "Datos de contacto, nombre y horario del negocio.",
			},
			"menu": {
				Type:        genai.TypeArray,
				Items:       menuCategorySchema,
				Description: "Lista de categorías que componen el menú principal.",
			},
		},
		Required: []string{"business", "menu"},
	}
}
