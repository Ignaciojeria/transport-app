package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"

	// Usar la librería oficial de Google para Vertex AI
	"google.golang.org/genai"
)

type AIOptimizeFleetRequestVehiclesExtractor func(input interface{}) (any, error)

func init() {
	ioc.Registry(NewAIOptimizeFleetRequestVehiclesExtractor, NewClient, httpserver.New)
}

func NewAIOptimizeFleetRequestVehiclesExtractor(client *genai.Client, s httpserver.Server) AIOptimizeFleetRequestVehiclesExtractor {
	return func(input interface{}) (any, error) {
		// Crear el esquema GenAI directamente para vehículos
		genaiVehiclesSchema := createVehiclesSchema()

		// Verifica que el esquema se haya generado correctamente
		if genaiVehiclesSchema == nil {
			return nil, fmt.Errorf("failed to transform OpenAPI schema to GenAI schema")
		}

		// Log del esquema generado para debug
		fmt.Printf("Esquema GenAI generado: %+v\n", genaiVehiclesSchema)

		// Usar la API correcta de google.golang.org/genai
		inputJSON, err := json.Marshal(input)
		if err != nil {
			return nil, fmt.Errorf("error serializando input: %w", err)
		}

		fmt.Printf("Input para extracción de vehículos: %s\n", string(inputJSON))

		prompt := fmt.Sprintf(`
You are a helpful assistant that extracts only the vehicles from a given input and returns them in a JSON format.

Your task is to analyze the following input and extract only the vehicle information, structuring it according to the provided schema.

Input to analyze:
%s

Instructions:
1. Extract only the vehicle information from the input.
2. If there are no vehicles in the input, return an empty array.
3. Ensure all required fields are present.
4. If a field is not available in the input, use empty values.
5. Respond with a JSON object containing a "vehicles" key with the array of vehicles.
6. Example format: {"vehicles": [...]}

CRITICAL DATA TYPE REQUIREMENTS:
- ALL numeric fields MUST be integers - this is MANDATORY for API compatibility
- Convert ALL string numbers to integers: "8000" becomes 8000, "12" becomes 12
- Volume, weight, insurance, serviceTime, and quantity MUST be integers, never strings
- If a numeric field is missing, use 0 as default integer value
- NEVER leave numeric fields as strings - always convert to integers

IMPORTANT PLATE EXTRACTION RULES:
- If you find a "plate" or "license_plate" or "patent" or "patente" field, use it as the vehicle's plate.
- If no plate field is found, look for an "id" field that could be a plate (alphanumeric codes, typically 6-8 characters).
- If you find an "id" field that looks like a vehicle identifier (alphanumeric, not purely numeric), consider it as the plate.
- Look for fields like "vehicle_id", "car_id", "plate", "registration", "number" that might contain patent information.
- If multiple candidates exist, prioritize fields that look more like a patent (alphanumeric, reasonable length).
- Always ensure the patent field is populated with the most appropriate identifier found.

IMPORTANT: Always return a JSON object with a "vehicles" key, even if the array is empty.
`, string(inputJSON))

		// Crear el contenido usando la API correcta
		content := genai.Text(prompt)

		ctx := context.Background()
		resp, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash", content, &genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema:   genaiVehiclesSchema,
		})
		if err != nil {
			return nil, fmt.Errorf("error generating content with Gemini: %w", err)
		}

		// Usar el método Text() para obtener la respuesta directamente
		output := resp.Text()

		// Log de la respuesta para debug
		fmt.Printf("Respuesta de Gemini: %s\n", output)

		// Intentar deserializar como objeto primero
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(output), &result); err != nil {
			// Si falla, intentar como array directamente
			var vehiclesArray []interface{}
			if err := json.Unmarshal([]byte(output), &vehiclesArray); err != nil {
				return nil, fmt.Errorf("error parsing Gemini response: %w", err)
			}
			// Si es un array, devolverlo directamente
			fmt.Printf("Respuesta procesada como array directo\n")
			return vehiclesArray, nil
		}

		// Si es un objeto, buscar la clave "vehicles"
		if vehicles, ok := result["vehicles"]; ok {
			fmt.Printf("Respuesta procesada como objeto con clave 'vehicles'\n")
			return vehicles, nil
		}

		// Si no hay clave "vehicles", devolver el objeto completo
		fmt.Printf("Respuesta procesada como objeto completo\n")
		return result, nil
	}
}

// createVehiclesSchema crea un esquema GenAI específico para vehículos.
func createVehiclesSchema() *genai.Schema {
	// Crear un esquema específico para vehículos basado en el schema real
	vehicleSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"plate": {
				Type:        genai.TypeString,
				Description: "Vehicle license plate or identifier",
			},
			"capacity": {
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"volume": {
						Type:        genai.TypeNumber,
						Description: "Volume capacity in cubic centimeters",
					},
					"weight": {
						Type:        genai.TypeNumber,
						Description: "Weight capacity in grams",
					},
					"deliveryUnitsQuantity": {
						Type:        genai.TypeInteger,
						Description: "Number of delivery units the vehicle can carry",
					},
					"insurance": {
						Type:        genai.TypeInteger,
						Description: "Insurance value in currency units",
					},
				},
			},
			"skills": {
				Type:        genai.TypeArray,
				Items:       &genai.Schema{Type: genai.TypeString},
				Description: "Skills required for the vehicle",
			},
		},
	}

	// Crear el esquema de respuesta final (array de vehículos)
	responseSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"vehicles": {
				Type:  genai.TypeArray,
				Items: vehicleSchema,
			},
		},
	}

	fmt.Printf("Esquema final generado para vehículos: %+v\n", responseSchema)
	return responseSchema
}
