package schema

import "google.golang.org/genai"

// GetImageGenerationRequestSchema define la estructura para una solicitud de generación de imagen.
func GetImageGenerationRequestSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"menuItemId": {
				Type:        genai.TypeString,
				Description: "ID del MenuItem, Side, o identificador especial para imágenes del menú. Valores especiales: 'footer' para la imagen del footer (footerImage). Para items o sides, debe corresponder al campo 'id' del elemento. NOTA: Para imágenes de portada, usa el campo 'coverImageGenerationRequest' en lugar de este array.",
			},
			"prompt": {
				Type:        genai.TypeString,
				Description: "Descripción profesional en inglés. OBLIGATORIO: usar EXACTAMENTE los ingredientes del array 'description' del item; no inventar ni omitir ninguno. Para sushi/piezas/rolls: incluir CADA variedad con su contenido Y envoltorio (Env) exactos. Ej: 'Professional food photography of [title]. Show exactly: [pieza 1: contenido + Env X], [pieza 2: contenido + Env Y]. Each piece must show correct filling and wrapper.'",
			},
			"aspectRatio": {
				Type:        genai.TypeString,
				Description: "Proporción de la imagen. Por defecto debe ser '1:1' para imágenes cuadradas. NOTA: Para imágenes de portada, usa el campo 'coverImageGenerationRequest' que automáticamente usa aspect ratio 3:1.",
			},
			"imageCount": {
				Type:        genai.TypeInteger,
				Description: "Cantidad de imágenes a generar. Por defecto debe ser 1.",
			},
		},
		Required: []string{"menuItemId", "prompt", "aspectRatio", "imageCount"},
	}
}
