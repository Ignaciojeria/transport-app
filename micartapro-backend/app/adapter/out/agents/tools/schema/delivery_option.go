package schema

import "google.golang.org/genai"

// GetTimeWindowSchema define la estructura para una ventana de tiempo.
func GetTimeWindowSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"start": {
				Type:        genai.TypeString,
				Description: "Hora de inicio de la ventana de tiempo (formato: HH:MM o ISO 8601).",
			},
			"end": {
				Type:        genai.TypeString,
				Description: "Hora de fin de la ventana de tiempo (formato: HH:MM o ISO 8601).",
			},
		},
		Required: []string{"start", "end"},
	}
}

// GetDeliveryOptionSchema define la estructura para una opción de envío/retiro.
func GetDeliveryOptionSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"type": {
				Type:        genai.TypeString,
				Enum:        []string{"PICKUP", "DELIVERY", "DIGITAL"},
				Description: "Tipo de opción: 'PICKUP' (retiro en tienda), 'DELIVERY' (envío a domicilio) o 'DIGITAL' (productos digitales, entrega instantánea).",
			},
			"requireTime": {
				Type:        genai.TypeBoolean,
				Description: "Si es true, el cliente debe indicar un tiempo para esta opción (hora exacta o ventana). Si es false, no se solicita tiempo.",
			},
			"timeRequestType": {
				Type:        genai.TypeString,
				Enum:        []string{"EXACT", "WINDOW"},
				Description: "Define el tipo de input que debe pedir la UI cuando requireTime=true: EXACT (hora exacta) o WINDOW (rango). Si requireTime=false, debe omitirse.",
			},
			"timeWindows": {
				Type:        genai.TypeArray,
				Items:       GetTimeWindowSchema(),
				Description: "Ventanas válidas/restricciones de atención para esta opción. Si se proveen, el tiempo elegido (EXACT o WINDOW) debe caer dentro de alguna ventana. Si se omite o está vacío, no hay restricción.",
			},
		},
		Required: []string{"type", "requireTime"},
	}
}
