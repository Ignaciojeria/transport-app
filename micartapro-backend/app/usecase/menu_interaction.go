package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"micartapro/app/adapter/out/agents"
	"micartapro/app/domain"
	"micartapro/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type MenuInteraction func(ctx context.Context, input domain.MenuInteractionRequest) (string, error)

func init() {
	ioc.Registry(
		NewMenuInteraction,
		observability.NewObservability,
		agents.NewMenuInteractionAgent)
}

func NewMenuInteraction(
	obs observability.Observability,
	menuInteractionAgent agents.MenuInteractionAgent) MenuInteraction {
	return func(ctx context.Context, input domain.MenuInteractionRequest) (string, error) {
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

			var createRequest domain.MenuCreateRequest

			// Mapeo directo del JSON crudo a tu objeto de dominio
			if err := json.Unmarshal(agentResp.CommandArgs, &createRequest); err != nil {
				return "", fmt.Errorf("error al mapear a MenuCreateRequest: %w", err)
			}

			// 3. Ejecutar Lógica de Negocio (Guardar, Emitir Evento, etc.)
			//menuService.ProcessCreateMenuRequest(ctx, createRequest, input.UserID)

			// 4. Retornar mensaje de éxito al usuario
			return "¡Menú creado exitosamente! Se ha notificado a todos los servicios.", nil
		}

		// ... manejar otros comandos

		return "Lo siento, no pude entender la acción solicitada.", nil
	}
}
