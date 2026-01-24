package schema

import "google.golang.org/genai"

// GetImageEditionRequestSchema define la estructura para una solicitud de edición de imagen.
func GetImageEditionRequestSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"menuItemId": {
				Type:        genai.TypeString,
				Description: "ID del MenuItem, Side, o identificador especial para imágenes del menú. Valores especiales: 'footer' para la imagen del footer (footerImage). Para items o sides, debe corresponder al campo 'id' del elemento.",
			},
			"prompt": {
				Type:        genai.TypeString,
				Description: "Descripción profesional y detallada en inglés para la edición o generación de la imagen basada en la imagen de referencia, enfocada en fotografía gastronómica profesional. Describe los cambios o mejoras que se deben aplicar a la imagen de referencia (ej. 'Add more vibrant colors and professional lighting to the food photography').",
			},
			"aspectRatio": {
				Type:        genai.TypeString,
				Description: "Proporción de la imagen. Por defecto debe ser '1:1' para imágenes cuadradas. Si no se especifica, se mantendrá el aspect ratio de la imagen de referencia.",
			},
			"imageCount": {
				Type:        genai.TypeInteger,
				Description: "Cantidad de imágenes a generar o editar. Por defecto debe ser 1.",
			},
			"referenceImageUrl": {
				Type:        genai.TypeString,
				Description: "URL completa de la imagen de referencia que será descargada y utilizada como base para la edición. La URL puede venir de DOS fuentes: 1) Del menú existente: usa la URL del campo 'photoUrl' del MenuItem o Side correspondiente en el menú actual, o del campo 'coverImage'/'footerImage' si es para imágenes especiales. 2) De una URL proporcionada por el usuario: si el usuario proporciona una URL específica en su solicitud o en el campo [FOTO_ADJUNTA], usa esa URL. La URL debe ser accesible y válida.",
			},
		},
		Required: []string{"menuItemId", "prompt", "referenceImageUrl", "aspectRatio", "imageCount"},
	}
}
