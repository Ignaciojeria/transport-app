package agents

import (
	"google.golang.org/genai"
)

// MenuInteractionResponse devuelve un schema para la estructura de respuesta de interacción con el menú
func MenuInteractionResponse() *genai.Schema {
	// Schema para un side (acompañamiento)
	sideSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"id":    {Type: genai.TypeString},
			"name":  {Type: genai.TypeString},
			"price": {Type: genai.TypeNumber},
		},
		Required: []string{"id", "name", "price"},
	}

	// Schema para un item del menú
	// Nota: Un item puede tener:
	// - title (requerido)
	// - description (opcional, puede estar vacío)
	// - sides (opcional, array de acompañamientos con precios)
	// - price (opcional, precio directo cuando no hay sides)
	menuItemSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"title":       {Type: genai.TypeString},
			"description": {Type: genai.TypeString},
			"sides": {
				Type:  genai.TypeArray,
				Items: sideSchema,
			},
			"price": {Type: genai.TypeNumber},
		},
		Required: []string{"title"},
	}

	// Schema para una categoría del menú
	menuCategorySchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"title": {Type: genai.TypeString},
			"items": {
				Type:  genai.TypeArray,
				Items: menuItemSchema,
			},
		},
		Required: []string{"title", "items"},
	}

	// Schema para contact
	contactSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"whatsapp": {Type: genai.TypeString},
		},
		Required: []string{"whatsapp"},
	}

	// Schema raíz del menú completo
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"contact": contactSchema,
			"businessHours": {
				Type:  genai.TypeArray,
				Items: &genai.Schema{Type: genai.TypeString},
			},
			"menu": {
				Type:  genai.TypeArray,
				Items: menuCategorySchema,
			},
		},
		Required: []string{"contact", "businessHours", "menu"},
	}
}
