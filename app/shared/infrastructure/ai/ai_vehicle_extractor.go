package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/getkin/kin-openapi/openapi3"

	// Usar la librería oficial de Google para Vertex AI
	"google.golang.org/genai"
)

type AIOptimizeFleetRequestVehiclesExtractor func(input interface{}) (any, error)

func init() {
	ioc.Registry(NewAIOptimizeFleetRequestVehiclesExtractor, NewClient, httpserver.New)
}

func NewAIOptimizeFleetRequestVehiclesExtractor(client *genai.Client, s httpserver.Server) AIOptimizeFleetRequestVehiclesExtractor {
	return func(input interface{}) (any, error) {
		fleetsSchema, ok := s.Manager.OpenAPI.Description().Components.Schemas["OptimizeFleetRequest"]

		if !ok {
			return nil, fmt.Errorf("optimize fleet request schema not found")
		}

		// Transforma el esquema de OpenAPI a un esquema de GenAI una sola vez
		genaiVehiclesSchema := createVehiclesSchemaFromOpenAPI(fleetsSchema.Value)

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

IMPORTANT PATENT EXTRACTION RULES:
- If you find a "patent" or "license_plate" field, use it as the vehicle's patent.
- If no patent field is found, look for an "id" field that could be a patent (alphanumeric codes, typically 6-8 characters).
- If you find an "id" field that looks like a vehicle identifier (alphanumeric, not purely numeric), consider it as the patent.
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

// createVehiclesSchemaFromOpenAPI crea un esquema GenAI basado en el esquema OpenAPI.
func createVehiclesSchemaFromOpenAPI(openAPISchema *openapi3.Schema) *genai.Schema {
	if openAPISchema == nil {
		fmt.Printf("OpenAPI schema es nil\n")
		return nil
	}

	fmt.Printf("Procesando esquema OpenAPI: %+v\n", openAPISchema)

	genaiSchema := &genai.Schema{}

	// Convertir el tipo de OpenAPI al tipo de GenAI
	// Maneja el caso en que el tipo es un slice de strings (multi-tipo)
	var openAPIType string
	// Corrected: Check if the pointer is not nil before dereferencing and getting the length
	if openAPISchema.Type != nil && len(*openAPISchema.Type) > 0 {
		openAPIType = (*openAPISchema.Type)[0]
	}

	fmt.Printf("Tipo OpenAPI detectado: %s\n", openAPIType)

	switch openAPIType {
	case "object":
		genaiSchema.Type = genai.TypeObject
		genaiSchema.Properties = make(map[string]*genai.Schema)
		for propName, propSchema := range openAPISchema.Properties {
			genaiSchema.Properties[propName] = createVehiclesSchemaFromOpenAPI(propSchema.Value)
		}
	case "array":
		genaiSchema.Type = genai.TypeArray
		if openAPISchema.Items != nil {
			genaiSchema.Items = createVehiclesSchemaFromOpenAPI(openAPISchema.Items.Value)
		}
	case "string":
		genaiSchema.Type = genai.TypeString
	case "number":
		genaiSchema.Type = genai.TypeNumber
	case "integer":
		genaiSchema.Type = genai.TypeInteger
	case "boolean":
		genaiSchema.Type = genai.TypeBoolean
	default:
		// Manejar tipos no soportados o desconocidos
		fmt.Printf("Unsupported OpenAPI type: %s\n", openAPIType)
		return nil
	}

	// Extraer y convertir el sub-esquema de "vehicles"
	if openAPISchema.Properties != nil {
		if vehiclesProp, ok := openAPISchema.Properties["vehicles"]; ok && vehiclesProp != nil {
			fmt.Printf("Encontrada propiedad 'vehicles' en el esquema\n")
			// El esquema de respuesta final debe ser el array de vehículos, no el objeto completo
			return createVehiclesSchemaFromOpenAPI(vehiclesProp.Value)
		}
	}

	// Considerar el caso donde el esquema del array está definido en el nivel raíz
	if genaiSchema.Type == genai.TypeArray && openAPISchema.Items != nil {
		fmt.Printf("Esquema de array detectado en nivel raíz\n")
		return genaiSchema
	}

	fmt.Printf("Esquema final generado: %+v\n", genaiSchema)
	return genaiSchema
}
