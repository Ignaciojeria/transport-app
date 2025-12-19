package schema

import "google.golang.org/genai"

// GetSideSchema define la estructura para un acompañamiento (ej. extra de queso).
func GetSideSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"id":    {Type: genai.TypeString, Description: "Identificador único (UUID) del acompañamiento."},
			"name":  {Type: genai.TypeString, Description: "Nombre del acompañamiento (ej. 'Extra de tocino')."},
			"price": {Type: genai.TypeNumber, Description: "Precio en número entero o flotante."},
		},
		Required: []string{"name", "price"},
	}
}
