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

				// 1. businessInfo: Agrupa los datos del negocio y reutiliza el Schema de Contacto completo
				"businessInfo": {
					Type:        genai.TypeObject,
					Description: "Información del negocio: nombre, contacto de WhatsApp y horario de atención, extraída del prompt del usuario.",
					Properties:  contactSchema.Properties,
					Required:    contactSchema.Required,
				},
				// 2. menuItemsStructure: Reemplaza el TypeString por la estructura de la carta (Array de Categorías)
				"menu": {
					Type:        genai.TypeArray,
					Description: "Estructura completa de categorías e ítems del menú. El modelo debe convertir el texto del usuario a este Array de objetos.",
					// ¡CLAVE! Usamos el Schema de Categoría para cada ítem del array.
					Items: menuCategorySchema,
				},
				// 3. coverImage: Imagen de portada del menú
				"coverImage": {
					Type:        genai.TypeString,
					Description: "URL de la imagen de portada del menú digital.",
				},
				// 4. footerImage: Imagen del footer/logo del menú
				"footerImage": {
					Type:        genai.TypeString,
					Description: "URL de la imagen del footer o logo del menú digital.",
				},
				// 5. deliveryOptions: Opciones de envío/retiro disponibles
				"deliveryOptions": {
					Type:        genai.TypeArray,
					Description: "Lista opcional de opciones de envío/retiro disponibles (PICKUP o DELIVERY).",
					Items:       schema.GetDeliveryOptionSchema(),
				},
			},
			// Los requeridos son las estructuras complejas que encapsulan todo
			//Required: []string{"businessInfo", "menu"},
		},
	}
}
