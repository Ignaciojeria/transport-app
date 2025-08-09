package agents

import (
	"context"
	"encoding/json"
	"fmt"

	"transport-app/app/adapter/out/agents/model"
	"transport-app/app/shared/infrastructure/ai"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"google.golang.org/genai"
)

type VehicleFieldNamesNormalizer func(ctx context.Context, input interface{}) (map[string]string, error)

func init() {
	ioc.Registry(NewVehicleFieldNamesNormalizer, ai.NewClient)
}

func NewVehicleFieldNamesNormalizer(client *genai.Client) VehicleFieldNamesNormalizer {
	vs := model.NewVehicleFieldMappingSchema()

	const modelName = "gemini-2.0-flash"

	return func(ctx context.Context, input interface{}) (map[string]string, error) {
		// 1) Prompt generado por VehicleFieldMappingSchema
		prompt := vs.Prompt(input)

		resp, err := client.Models.GenerateContent(
			ctx,
			modelName,
			genai.Text(prompt),
			&genai.GenerateContentConfig{
				ResponseMIMEType: "application/json",
				ResponseSchema:   model.SchemaForVehicleFieldMapping(),
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
