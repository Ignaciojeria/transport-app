package schema

import "google.golang.org/genai"

// GetCoverImageGenerationRequestSchema define la estructura para una solicitud de generación de imagen de portada.
func GetCoverImageGenerationRequestSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"prompt": {
				Type:        genai.TypeString,
				Description: "Descripción profesional y detallada en inglés para la generación de la imagen de portada, enfocada en crear una imagen visual atractiva que represente el estilo del menú o negocio. Debe reflejar la identidad visual del restaurante o negocio (ej. 'Modern restaurant cover image with elegant food presentation, warm lighting, professional photography style'). La imagen será horizontalmente larga y verticalmente corta, tipo foto portada LinkedIn (aspect ratio 16:9).",
			},
			"imageCount": {
				Type:        genai.TypeInteger,
				Description: "Cantidad de imágenes a generar. Por defecto debe ser 1.",
			},
		},
		Required: []string{"prompt", "imageCount"},
	}
}
