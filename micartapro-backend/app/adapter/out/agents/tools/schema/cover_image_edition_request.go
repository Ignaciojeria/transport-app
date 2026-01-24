package schema

import "google.golang.org/genai"

// GetCoverImageEditionRequestSchema define la estructura para una solicitud de edición de imagen de portada.
func GetCoverImageEditionRequestSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"prompt": {
				Type:        genai.TypeString,
				Description: "Descripción profesional y detallada en inglés para la edición de la imagen de portada basada en la imagen de referencia, enfocada en crear una imagen visual atractiva que represente el estilo del menú o negocio. Describe los cambios o mejoras que se deben aplicar a la imagen de referencia (ej. 'Add more vibrant colors, enhance the lighting, and improve the professional photography style while maintaining the restaurant identity'). La imagen será horizontalmente larga y verticalmente corta, tipo foto portada LinkedIn (aspect ratio 16:9).",
			},
			"imageCount": {
				Type:        genai.TypeInteger,
				Description: "Cantidad de imágenes a generar o editar. Por defecto debe ser 1.",
			},
			"referenceImageUrl": {
				Type:        genai.TypeString,
				Description: "URL completa de la imagen de portada de referencia que será descargada y utilizada como base para la edición. La URL puede venir de DOS fuentes: 1) Del menú existente: usa la URL del campo 'coverImage' del menú actual si el usuario quiere editar la imagen de portada existente. 2) De una URL proporcionada por el usuario: si el usuario proporciona una URL específica en su solicitud o en el campo [FOTO_ADJUNTA], usa esa URL. La URL debe ser accesible y válida.",
			},
		},
		Required: []string{"prompt", "imageCount", "referenceImageUrl"},
	}
}
