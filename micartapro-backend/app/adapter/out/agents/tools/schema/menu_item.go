package schema

import "google.golang.org/genai"

// GetMenuItemSchema define la estructura para un producto individual del menú.
func GetMenuItemSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"id": {
				Type:        genai.TypeString,
				Description: "Identificador semántico único del item en formato kebab-case (ej. 'empanadas-pino', 'pizza-margherita').",
			},
			"title": {
				Type:        genai.TypeObject,
				Description: "Título del producto en formato multiidioma con base y traducciones.",
				Properties: map[string]*genai.Schema{
					"base": {
						Type:        genai.TypeString,
						Description: "Texto base del título (idioma principal, generalmente español).",
					},
					"languages": {
						Type:        genai.TypeObject,
						Description: "Traducciones del título en diferentes idiomas.",
						Properties: map[string]*genai.Schema{
							"es": {Type: genai.TypeString, Description: "Título en español."},
							"en": {Type: genai.TypeString, Description: "Título en inglés."},
							"pt": {Type: genai.TypeString, Description: "Título en portugués."},
						},
					},
				},
				Required: []string{"base", "languages"},
			},
			"description": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"id": {
							Type:        genai.TypeString,
							Description: "ID único del bloque (kebab-case). OBLIGATORIO si tiene selectables. Ej: 'california-camaron-palta'. Permite mapear descriptionSelections en el pedido.",
						},
						"base": {
							Type:        genai.TypeString,
							Description: "Texto base de este elemento (idioma principal, generalmente español).",
						},
						"languages": {
							Type:        genai.TypeObject,
							Description: "Traducciones de este elemento en diferentes idiomas.",
							Properties: map[string]*genai.Schema{
								"es": {Type: genai.TypeString, Description: "Texto en español."},
								"en": {Type: genai.TypeString, Description: "Texto en inglés."},
								"pt": {Type: genai.TypeString, Description: "Texto en portugués."},
							},
						},
						"selectables": {
							Type:        genai.TypeObject,
							Description: "Preferencias sin precio dentro de este bloque (ej. 'envuelto en sésamo o ciboulette'). Solo cuando el texto indica opciones elegibles que NO cambian precio.",
							Properties: map[string]*genai.Schema{
								"selection": {
									Type: genai.TypeObject,
									Properties: map[string]*genai.Schema{
										"mode": {Type: genai.TypeString, Enum: []string{"SINGLE", "MULTIPLE"}, Description: "SINGLE: elegir una. MULTIPLE: elegir varias."},
										"min":  {Type: genai.TypeInteger, Description: "Mínimo a elegir (0 = opcional)."},
										"max":  {Type: genai.TypeInteger, Description: "Máximo a elegir."},
									},
								},
								"options": {
									Type: genai.TypeArray,
									Items: &genai.Schema{
										Type: genai.TypeObject,
										Properties: map[string]*genai.Schema{
											"id":   {Type: genai.TypeString, Description: "ID único de la opción (kebab-case)."},
											"name": {Type: genai.TypeObject, Description: "Nombre en formato multiidioma.", Properties: map[string]*genai.Schema{"base": {Type: genai.TypeString}, "languages": {Type: genai.TypeObject, Properties: map[string]*genai.Schema{"es": {Type: genai.TypeString}, "en": {Type: genai.TypeString}, "pt": {Type: genai.TypeString}}}}},
										},
										Required: []string{"id", "name"},
									},
								},
							},
						},
					},
					Required: []string{"base", "languages"},
				},
				Description: "Array de descripciones: cada elemento = una dimensión (ingredientes, preparación, notas). Con id y selectables opcionales para preferencias sin precio.",
			},
			"foodAttributes": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeString,
					Enum: []string{"GLUTEN", "SEAFOOD", "NUTS", "DAIRY", "EGGS", "SOY", "VEGAN", "VEGETARIAN", "SPICY", "ALCOHOL"},
				},
				Description: "Atributos alimentarios opcionales: alérgenos (GLUTEN, SEAFOOD, NUTS, DAIRY, EGGS, SOY), dieta (VEGAN, VEGETARIAN), SPICY, ALCOHOL.",
			},
			"sides": {
				Type: genai.TypeArray,
				// ¡Aquí se utiliza el Schema del Side!
				Items:       GetSideSchema(),
				Description: "Lista opcional de acompañamientos o variaciones con pricing (ej. tamaño grande, extra de salsa).",
			},
			"pricing": {
				Type:        genai.TypeObject,
				Description: "Estructura de pricing del ítem del menú.",
				Properties:  GetPricingSchema().Properties,
				Required:    GetPricingSchema().Required,
			},
			"photoUrl": {
				Type:        genai.TypeString,
				Description: "URL opcional de la imagen del producto. Debe ser una URL pública accesible.",
			},
			"station": {
				Type:        genai.TypeString,
				Enum:        []string{"KITCHEN", "BAR"},
				Description: "Estación que prepara este ítem: KITCHEN (cocina) o BAR (bar). Opcional.",
			},
		},
		Required: []string{"id", "title", "pricing"},
	}
}
