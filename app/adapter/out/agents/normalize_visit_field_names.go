package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"transport-app/app/adapter/out/agents/model"
	"transport-app/app/shared/infrastructure/ai"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"google.golang.org/genai"
)

type VisitFieldNamesNormalizer func(input interface{}) (map[string]string, error)

func init() {
	ioc.Registry(NewVisitFieldNamesNormalizer, ai.NewClient)
}

func NewVisitFieldNamesNormalizer(client *genai.Client) VisitFieldNamesNormalizer {
	vs := model.NewVisitFieldMappingSchema()

	const modelName = "gemini-2.0-flash"

	return func(input interface{}) (map[string]string, error) {
		// 1) Prompt generado por VisitFieldMappingSchema
		prompt := vs.Prompt(input)

		// 2) Llamada a Gemini con schema para mapeo de claves
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		resp, err := client.Models.GenerateContent(
			ctx,
			modelName,
			genai.Text(prompt),
			&genai.GenerateContentConfig{
				ResponseMIMEType: "application/json",
				ResponseSchema:   model.SchemaForFieldMapping(),
			},
		)
		if err != nil {
			return nil, fmt.Errorf("genai error: %w", err)
		}

		out := resp.Text()

		// 3) Unmarshal directo a map[string]string
		var keyMapping map[string]string
		if err := json.Unmarshal([]byte(out), &keyMapping); err != nil {
			return nil, fmt.Errorf("no pude parsear la respuesta LLM: %w", err)
		}

		// 4) Filtrar campos vacíos
		result := make(map[string]string)
		for originalKey, officialKey := range keyMapping {
			if originalKey != "" {
				result[originalKey] = officialKey
			}
		}

		return result, nil
	}
}
