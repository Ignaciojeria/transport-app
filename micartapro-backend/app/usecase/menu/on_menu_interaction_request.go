package menu

import (
	"context"
	"encoding/json"
	"fmt"
	"micartapro/app/adapter/out/agents"
	"micartapro/app/adapter/out/storage"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/eventprocessing"
	"micartapro/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type OnMenuInteractionRequest func(ctx context.Context, input events.MenuInteractionRequest) (string, error)

func init() {
	ioc.Registry(
		NewOnMenuInteractionRequest,
		observability.NewObservability,
		agents.NewMenuInteractionAgent,
		eventprocessing.NewPublisherStrategy,
		storage.NewGetLatestMenuById)
}

func NewOnMenuInteractionRequest(
	obs observability.Observability,
	menuInteractionAgent agents.MenuInteractionAgent,
	publisherManager eventprocessing.PublisherManager,
	getLatestMenuById storage.GetLatestMenuById) OnMenuInteractionRequest {
	return func(ctx context.Context, input events.MenuInteractionRequest) (string, error) {

		menu, err := getLatestMenuById(ctx, input.MenuID)
		if err != nil && err != storage.ErrMenuNotFound {
			return "", err
		}
		input.JsonMenu = menu

		// 1. Llamar al Agente Procesador (que hace la inspección del FunctionCall)
		agentResp, err := menuInteractionAgent(ctx, input)
		if err != nil {
			return "", err
		}

		// 2. Lógica Limpia: ¿Es texto o un comando?

		if agentResp.TextResponse != "" {
			// Caso 1: Texto conversacional (Ej. "Faltan los precios")
			return agentResp.TextResponse, nil
		}

		if agentResp.CommandName == "createMenu" {
			// Caso 2: Comando. Hacemos el mapeo (la única deserialización necesaria)
			var createRequest events.MenuCreateRequest

			// Mapeo directo del JSON crudo a tu objeto de dominio
			if err := json.Unmarshal(agentResp.CommandArgs, &createRequest); err != nil {
				return "", fmt.Errorf("error al mapear a MenuCreateRequest: %w", err)
			}
			createRequest.ID = input.MenuID

			if err := publisherManager.Publish(ctx, eventprocessing.PublishRequest{
				Topic:  "micartapro.events",
				Source: "micartapro.agent.menu.interaction",
				Event:  createRequest,
			}); err != nil {
				return "", fmt.Errorf("error al publicar evento de creación de menú: %w", err)
			}

			return "¡Menú creado exitosamente! Se ha notificado a todos los servicios.", nil
		}

		// ... manejar otros comandos

		return "Lo siento, no pude entender la acción solicitada.", nil
	}
}
