package agents

import (
	"google.golang.org/genai"
)

func MenuInteractionResponse() *genai.Schema {

	// --- UnitOfMeasure enum ---
	unitSchema := &genai.Schema{
		Type: genai.TypeString,
		Enum: []string{
			"EACH",
			"GRAM",
			"KILOGRAM",
			"MILLILITER",
			"LITER",
		},
	}

	// --- PricingMode enum ---
	modeSchema := &genai.Schema{
		Type: genai.TypeString,
		Enum: []string{
			"UNIT",
			"WEIGHT",
			"VOLUME",
		},
	}

	// --- Pricing schema (currency es global en businessInfo, no en pricing) ---
	pricingSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"mode":         modeSchema,
			"unit":         unitSchema,
			"pricePerUnit": {Type: genai.TypeNumber},
			"baseUnit":     {Type: genai.TypeNumber},
		},
		Required: []string{"mode", "unit", "pricePerUnit"},
	}

	// --- Side ---
	sideSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"id":      {Type: genai.TypeString},
			"name":    {Type: genai.TypeString},
			"pricing": pricingSchema,
			"station": {Type: genai.TypeString, Enum: []string{"KITCHEN", "BAR"}},
		},
		Required: []string{"id", "name", "pricing"},
	}

	// --- MenuItem ---
	menuItemSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"title": {Type: genai.TypeString},
			"description": {
				Type:  genai.TypeArray,
				Items: &genai.Schema{Type: genai.TypeString},
			},
			"sides": {
				Type:  genai.TypeArray,
				Items: sideSchema,
			},
			"pricing": pricingSchema,
			"station": {Type: genai.TypeString, Enum: []string{"KITCHEN", "BAR"}},
		},
		Required: []string{"title", "pricing"},
	}

	// --- MenuCategory ---
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

	// --- Contact (incluye currency global del negocio) ---
	contactSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"whatsapp": {Type: genai.TypeString},
			"currency": {Type: genai.TypeString, Enum: []string{"USD", "CLP", "BRL"}},
		},
		Required: []string{"whatsapp", "currency"},
	}

	// --- Root ---
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
