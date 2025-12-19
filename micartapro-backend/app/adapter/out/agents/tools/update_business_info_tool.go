package tools

import "google.golang.org/genai"

func UpdateBusinessInfoTool() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        "updateBusinessInfo",
		Description: "Actualiza información específica del negocio, como el número de WhatsApp, el nombre o los horarios de atención, sin modificar los ítems del menú.",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"businessName": {Type: genai.TypeString, Description: "El nuevo nombre del negocio."},
				"whatsapp":     {Type: genai.TypeString, Description: "El nuevo número de contacto de WhatsApp para pedidos."},
				"businessHours": {
					Type:        genai.TypeArray,
					Items:       &genai.Schema{Type: genai.TypeString},
					Description: "Los nuevos horarios de atención del negocio.",
				},
			},
			// IMPORTANTE: Ningún campo es Requerido, ya que el usuario puede querer actualizar solo uno (ej. solo WhatsApp).
			// Gemini extraerá solo el campo que el usuario mencione en el prompt.
			Required: nil,
		},
	}
}
