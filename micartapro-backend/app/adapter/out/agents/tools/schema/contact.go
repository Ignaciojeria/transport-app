package schema

import "google.golang.org/genai"

// GetContactSchema define la estructura para los datos de contacto y negocio.
func GetContactSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"businessName": {Type: genai.TypeString, Description: "El nombre oficial del negocio o restaurant."},
			"whatsapp":     {Type: genai.TypeString, Description: "Número de teléfono de contacto, incluyendo código de país (ej. +56912345678)."},

			// Mantenemos TypeArray, pero con una descripción mucho más estricta:
			"businessHours": {
				Type:  genai.TypeArray,
				Items: &genai.Schema{Type: genai.TypeString},
				// **DESCRIPCIÓN MEJORADA:** Instruye al modelo a usar un array y a separar los rangos.
				Description: "Lista de strings, donde cada string representa un rango de horario distinto. Debe ser un Array JSON (ej. ['Lunes a Viernes: 10h-22h', 'Sábado y Domingo: Cerrado']).",
			},
		},
		Required: []string{"whatsapp"},
	}
}
