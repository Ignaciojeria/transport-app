package schema

import "google.golang.org/genai"

// GetImageGenerationRequestSchema define la estructura para una solicitud de generación de imagen.
func GetImageGenerationRequestSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"menuItemId": {
				Type:        genai.TypeString,
				Description: "ID del MenuItem o Side que requiere la imagen. Debe corresponder al campo 'id' del elemento.",
			},
			"prompt": {
				Type:        genai.TypeString,
				Description: "Descripción profesional y detallada en inglés para la generación de la imagen, enfocada en fotografía gastronómica profesional (ej. 'Professional food photography of Chilean empanadas de pino on a wooden table').",
			},
			"aspectRatio": {
				Type:        genai.TypeString,
				Description: "Proporción de la imagen. Por defecto debe ser '1:1' para imágenes cuadradas.",
			},
			"imageCount": {
				Type:        genai.TypeInteger,
				Description: "Cantidad de imágenes a generar. Por defecto debe ser 1.",
			},
		},
		Required: []string{"menuItemId", "prompt", "aspectRatio", "imageCount"},
	}
}
