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
				// 6. coverImageGenerationRequest: Solicitud de generación de imagen de portada
				"coverImageGenerationRequest": {
					Type:        genai.TypeObject,
					Description: "Solicitud opcional de generación de imagen de portada. Solo incluir cuando el usuario solicita explícitamente generar o cambiar la imagen de portada. La imagen se generará con aspect ratio 16:9 (horizontalmente larga y verticalmente corta, tipo foto portada LinkedIn).",
					Properties:  schema.GetCoverImageGenerationRequestSchema().Properties,
					Required:    schema.GetCoverImageGenerationRequestSchema().Required,
				},
				// 7. coverImageEditionRequest: Solicitud de edición de imagen de portada
				"coverImageEditionRequest": {
					Type:        genai.TypeObject,
					Description: "Solicitud opcional de edición de imagen de portada. Solo incluir cuando el usuario solicita explícitamente editar, mejorar o modificar la imagen de portada existente. La URL de referencia puede venir del campo 'coverImage' del menú actual o de una URL proporcionada por el usuario en su solicitud o en [FOTO_ADJUNTA]. El agente debe proporcionar la URL completa de la imagen de referencia que se utilizará como base para la edición. La imagen será horizontalmente larga y verticalmente corta, tipo foto portada LinkedIn (aspect ratio 16:9).",
					Properties:  schema.GetCoverImageEditionRequestSchema().Properties,
					Required:    schema.GetCoverImageEditionRequestSchema().Required,
				},
				// 8. imageGenerationRequests: Solicitudes de generación de imágenes
				"imageGenerationRequests": {
					Type:        genai.TypeArray,
					Description: "Lista opcional de solicitudes de generación de imágenes para items o sides que requieren imagen. Solo incluir cuando el usuario solicita explícitamente generar o cambiar una imagen.",
					Items:       schema.GetImageGenerationRequestSchema(),
				},
				// 9. imageEditionRequests: Solicitudes de edición de imágenes
				"imageEditionRequests": {
					Type:        genai.TypeArray,
					Description: "Lista opcional de solicitudes de edición de imágenes para items o sides que requieren modificar una imagen existente. Solo incluir cuando el usuario solicita explícitamente editar, mejorar o modificar una imagen existente. La URL de referencia puede venir del campo 'photoUrl' del elemento correspondiente en el menú actual, del campo 'coverImage'/'footerImage' si es para imágenes especiales, o de una URL proporcionada por el usuario en su solicitud o en [FOTO_ADJUNTA]. El agente debe proporcionar la URL completa de la imagen de referencia que se utilizará como base para la edición.",
					Items:       schema.GetImageEditionRequestSchema(),
				},
			},
			// Los requeridos son las estructuras complejas que encapsulan todo
			//Required: []string{"businessInfo", "menu"},
		},
	}
}
