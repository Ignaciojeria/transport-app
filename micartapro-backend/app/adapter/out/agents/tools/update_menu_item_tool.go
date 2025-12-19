package tools

import "google.golang.org/genai"

func UpdateMenuItemTool() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        "updateMenuItem",
		Description: "Actualiza el precio, la descripción, o la categoría de un ítem existente en el menú.",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				// CLAVE: Saber qué ítem se está editando
				"itemName": {
					Type:        genai.TypeString,
					Description: "El nombre exacto del ítem que se desea modificar.",
				},
				// Lo que se actualizará (todos opcionales)
				"newPrice": {
					Type:        genai.TypeNumber,
					Description: "El nuevo precio del ítem.",
				},
				"newDescription": {
					Type:        genai.TypeString,
					Description: "La nueva descripción del ítem.",
				},
				"newCategory": {
					Type:        genai.TypeString,
					Description: "La nueva categoría a la que pertenecerá el ítem (ej. 'Platos principales').",
				},
				// Nota: No se incluye 'newName' por ser una operación riesgosa. Es mejor que el usuario
				// lo elimine y lo cree de nuevo si cambia radicalmente el nombre.
			},
			// Requerido saber el ítem a actualizar.
			Required: []string{"itemName"},
		},
	}
}
