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

type AIOptimizeFleetRequestVisitsExtractor func(input interface{}) (any, error)

func init() {
	ioc.Registry(NewAIOptimizeFleetRequestVisitsExtractor, NewClient, httpserver.New)
}

func NewAIOptimizeFleetRequestVisitsExtractor(client *genai.Client, s httpserver.Server) AIOptimizeFleetRequestVisitsExtractor {
	return func(input interface{}) (any, error) {
		// Crear el esquema GenAI directamente para visitas
		genaiVisitsSchema := createVisitsSchema()

		// Verifica que el esquema se haya generado correctamente
		if genaiVisitsSchema == nil {
			return nil, fmt.Errorf("failed to transform OpenAPI schema to GenAI schema")
		}

		// Log del esquema generado para debug
		fmt.Printf("Esquema GenAI generado para visits: %+v\n", genaiVisitsSchema)

		// Usar la API correcta de google.golang.org/genai
		inputJSON, err := json.Marshal(input)
		if err != nil {
			return nil, fmt.Errorf("error serializando input: %w", err)
		}

		fmt.Printf("Input para extracción de visits: %s\n", string(inputJSON))

		prompt := fmt.Sprintf(`
Eres un asistente de extracción de datos. Tu tarea es extraer la información de las visitas del JSON de entrada y devolverla como un objeto JSON que se ajuste exactamente al esquema proporcionado.

Entrada a analizar:
%s

Instrucciones:
- Extrae la información de las visitas, incluyendo todos los detalles del cliente y las coordenadas.
- Omite cualquier dato relacionado con 'pickup'.
- El output debe ser solo un JSON válido, sin explicaciones adicionales.
`, string(inputJSON))

		// Crear el contenido usando la API correcta
		content := genai.Text(prompt)

		ctx := context.Background()
		resp, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash", content, &genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema:   genaiVisitsSchema,
		})
		if err != nil {
			return nil, fmt.Errorf("error generating content with Gemini: %w", err)
		}

		// Usar el método Text() para obtener la respuesta directamente
		output := resp.Text()

		// Log de la respuesta para debug
		fmt.Printf("Respuesta de Gemini para visits: %s\n", output)

		// Intentar deserializar como objeto primero
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(output), &result); err != nil {
			// Si falla, intentar como array directamente
			var visitsArray []interface{}
			if err := json.Unmarshal([]byte(output), &visitsArray); err != nil {
				return nil, fmt.Errorf("error parsing Gemini response: %w", err)
			}
			// Si es un array, devolverlo directamente
			fmt.Printf("Respuesta procesada como array directo\n")
			return visitsArray, nil
		}

		// Si es un objeto, buscar la clave "visits"
		if visits, ok := result["visits"]; ok {
			fmt.Printf("Respuesta procesada como objeto con clave 'visits'\n")
			return visits, nil
		}

		// Si no hay clave "visits", devolver el objeto completo
		fmt.Printf("Respuesta procesada como objeto completo\n")
		return result, nil
	}
}

// createVisitsSchema crea un esquema GenAI específico para visitas.
func createVisitsSchema() *genai.Schema {
	// Definición de los campos de la visita según la estructura real del JSON
	visitSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"delivery": {
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"addressInfo": {
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"contact": {
								Type: genai.TypeObject,
								Properties: map[string]*genai.Schema{
									"fullName": {
										Type: genai.TypeString,
									},
									"phone": {
										Type: genai.TypeString,
									},
								},
							},
							"coordinates": {
								Type: genai.TypeObject,
								Properties: map[string]*genai.Schema{
									"latitude": {
										Type: genai.TypeNumber,
									},
									"longitude": {
										Type: genai.TypeNumber,
									},
								},
							},
						},
					},
				},
			},
			"orders": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"referenceID": {
							Type: genai.TypeString,
						},
						"deliveryUnits": {
							Type: genai.TypeArray,
							Items: &genai.Schema{
								Type: genai.TypeObject,
								Properties: map[string]*genai.Schema{
									"insurance": {
										Type: genai.TypeInteger,
									},
									"volume": {
										Type: genai.TypeInteger,
									},
									"weight": {
										Type: genai.TypeInteger,
									},
									"lpn": {
										Type: genai.TypeString,
									},
								},
							},
							Required: []string{"lpn"},
						},
					},
					Required: []string{"referenceID"},
				},
			},
		},
	}

	// El esquema de respuesta final debe ser un objeto con una clave 'visits' que contenga un array
	responseSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"visits": {
				Type:  genai.TypeArray,
				Items: visitSchema,
			},
		},
	}

	return responseSchema
}
