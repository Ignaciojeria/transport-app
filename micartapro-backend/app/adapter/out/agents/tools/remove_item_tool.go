package tools

import (
	"google.golang.org/genai"
	// Aunque no lo usamos directamente aquí, mantenemos la consistencia
	// import "micartapro-backend/app/adapter/out/agents/tools/schema"
)

func RemoveMenuItemTool() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        "removeMenuItem",
		Description: "Elimina permanentemente uno o varios ítems específicos del menú. Úsalo solo para ítems que existan en el MENU_ACTUAL.",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{

				"itemsToRemove": {
					Type:        genai.TypeArray,
					Description: "Lista de objetos, donde cada objeto identifica un ítem a eliminar.",
					Items: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{

							// 1. itemTitle (Requerido)
							"itemTitle": {
								Type:        genai.TypeString,
								Description: "El título exacto del ítem a eliminar (ej. 'Taco Mexicano').",
							},

							// 2. categoryTitle (Opcional, para evitar ambigüedad)
							"categoryTitle": {
								Type:        genai.TypeString,
								Description: "El título de la categoría donde se encuentra el ítem (ej. 'Platos Fuertes'). Opcional si el nombre del ítem es único.",
							},
						},
						// Solo el título del ítem es requerido
						Required: []string{"itemTitle"},
					},
				},
			},
			Required: []string{"itemsToRemove"},
		},
	}
}
