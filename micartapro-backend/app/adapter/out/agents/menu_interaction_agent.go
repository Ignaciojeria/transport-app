package agents

import (
	"context"
	"encoding/json"
	"fmt"

	"micartapro/app/adapter/out/agents/prompt"
	"micartapro/app/adapter/out/agents/tools"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/ai"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"google.golang.org/genai"
)

type MenuInteractionAgent func(ctx context.Context, req events.MenuInteractionRequest) (AgentResponse, error)

func init() {
	ioc.Registry(NewMenuInteractionAgent, ai.NewClient)
}

func NewMenuInteractionAgent(client *genai.Client) MenuInteractionAgent {
	const modelName = "gemini-2.5-pro"

	menuTools := tools.GetAllMenuTools()

	return func(ctx context.Context, req events.MenuInteractionRequest) (AgentResponse, error) {

		// 1. Construir el Prompt Completo y el Historial
		menuJSON := req.MenuToon()
		fullPrompt := prompt.MenuInteractionPrompt(menuJSON, req.Message, req.PhotoUrl)

		var contents []*genai.Content
		for _, msg := range req.History {
			contents = append(contents, &genai.Content{
				Role:  msg.Role,
				Parts: []*genai.Part{{Text: msg.Content}},
			})
		}
		contents = append(contents, &genai.Content{
			Role:  "user",
			Parts: []*genai.Part{{Text: fullPrompt}},
		})

		// 3. Llamada a GenerateContent con las Tools
		resp, err := client.Models.GenerateContent(ctx, modelName, contents, &genai.GenerateContentConfig{Tools: menuTools})
		if err != nil {
			return AgentResponse{}, fmt.Errorf("genai error: %w", err)
		}

		// 4. Procesar la Respuesta: Validaciones de seguridad
		if len(resp.Candidates) == 0 {
			return AgentResponse{}, fmt.Errorf("respuesta de Gemini vacía (sin candidatos)")
		}

		candidate := resp.Candidates[0]

		// **VALIDACIÓN IMPORTANTE**: Detectar si el modelo falló al armar la función
		if candidate.FinishReason == "MALFORMED_FUNCTION_CALL" {
			return AgentResponse{}, fmt.Errorf("error del modelo: llamada de función malformada (FinishReason: MALFORMED_FUNCTION_CALL)")
		}

		// Validar si Content o Parts son nulos antes de acceder al índice [0]
		if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
			return AgentResponse{}, fmt.Errorf("la respuesta no contiene partes válidas (FinishReason: %s)", candidate.FinishReason)
		}

		part := candidate.Content.Parts[0]

		if part.FunctionCall != nil {
			// **CASO 1: FUNCIÓN LLAMADA**
			// Serializamos *solo* los argumentos a json.RawMessage para retornarlo
			argsJSON, err := json.Marshal(part.FunctionCall.Args)
			if err != nil {
				return AgentResponse{}, fmt.Errorf("error al serializar FunctionCall args: %w", err)
			}

			return AgentResponse{
				CommandName: part.FunctionCall.Name,
				CommandArgs: argsJSON,
			}, nil
		}

		if part.Text != "" {
			// **CASO 2: TEXTO CONVERSACIONAL**
			return AgentResponse{
				TextResponse: part.Text,
			}, nil
		}

		// Caso de error inesperado o sin contenido
		return AgentResponse{}, fmt.Errorf("respuesta de Gemini no pudo ser interpretada como texto o FunctionCall")
	}
}
