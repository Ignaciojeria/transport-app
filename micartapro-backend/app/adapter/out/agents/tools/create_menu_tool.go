package tools

import (
	"google.golang.org/genai"
	// Asumimos que este es el path correcto donde están tus definiciones granulares de Schema
	"micartapro/app/adapter/out/agents/tools/schema"
)

// Definición de la Tool para crear el menú
func CreateMenuTool() *genai.FunctionDeclaration {

	// Obtenemos los schemas que definen las estructuras complejas
	contactSchema := schema.GetContactSchema()
	menuCategorySchema := schema.GetMenuCategorySchema()

	return &genai.FunctionDeclaration{
		// Nombre de la función en inglés para consistencia técnica
		Name: "createMenu",
		// Descripción en español para que el modelo lo relacione con los prompts del usuario
		Description: "Crea o actualiza el menú digital completo con los productos, precios, horario y datos de contacto.",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{

				// 1. businessInfo: Agrupa los datos del negocio y reutiliza el Schema de Contacto
				"businessInfo": {
					Type:        genai.TypeObject,
					Description: "Información del negocio: nombre, contacto de WhatsApp y horario de atención, extraída del prompt del usuario.",
					Properties: map[string]*genai.Schema{
						// businessName (antes era un argumento individual)
						"businessName": {Type: genai.TypeString, Description: "El nombre oficial del negocio."},
						// whatsapp (reutilizamos la definición del Schema de Contacto)
						"whatsapp": contactSchema.Properties["whatsapp"],
						// businessHours (reutilizamos la definición del Schema de Contacto)
						"businessHours": contactSchema.Properties["businessHours"],
					},
					//Required: []string{"businessName", "whatsapp"}, // businessHours puede ser opcional si el usuario no lo menciona
				},
				// 2. menuItemsStructure: Reemplaza el TypeString por la estructura de la carta (Array de Categorías)
				"menu": {
					Type:        genai.TypeArray,
					Description: "Estructura completa de categorías e ítems del menú. El modelo debe convertir el texto del usuario a este Array de objetos.",
					// ¡CLAVE! Usamos el Schema de Categoría para cada ítem del array.
					Items: menuCategorySchema,
				},
			},
			// Los requeridos son las estructuras complejas que encapsulan todo
			//Required: []string{"businessInfo", "menu"},
		},
	}
}
