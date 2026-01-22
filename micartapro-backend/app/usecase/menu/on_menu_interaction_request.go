package menu

import (
	"context"
	"encoding/json"
	"fmt"
	"micartapro/app/adapter/out/agents"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/eventprocessing"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type OnMenuInteractionRequest func(ctx context.Context, input events.MenuInteractionRequest) (string, error)

func init() {
	ioc.Registry(
		NewOnMenuInteractionRequest,
		observability.NewObservability,
		agents.NewMenuInteractionAgent,
		eventprocessing.NewPublisherStrategy,
		supabaserepo.NewGetMenuById)
}

func NewOnMenuInteractionRequest(
	obs observability.Observability,
	menuInteractionAgent agents.MenuInteractionAgent,
	publisherManager eventprocessing.PublisherManager,
	getMenuById supabaserepo.GetMenuById) OnMenuInteractionRequest {
	return func(ctx context.Context, input events.MenuInteractionRequest) (string, error) {

		// 1. Verificar y propagar el versionID del contexto para guardar la nueva versión generada
		// El versionID se usa cuando se guarda el menú, NO para leer el menú actual
		if versionID, ok := sharedcontext.VersionIDFromContext(ctx); ok && versionID != "" {
			obs.Logger.InfoContext(ctx, "versionID from context will be used when saving new menu version", "menuID", input.MenuID, "versionID", versionID)
		} else {
			obs.Logger.InfoContext(ctx, "no versionID in context, a new one will be generated when saving", "menuID", input.MenuID)
		}

		// 2. Obtener el contenido del menú desde Supabase usando el menuID
		// IMPORTANTE: Pasamos string vacío para que getMenuById SIEMPRE use current_version_id
		// NO usamos el versionID del contexto para leer, solo para guardar después
		obs.Logger.InfoContext(ctx, "getting menu content using current_version_id from menu (ignoring versionID from context)", "menuID", input.MenuID)

		// 3. Obtener el contenido del menú desde menu_versions en Supabase
		// Esta función obtiene el current_version_id del menú desde la tabla menus
		// y luego busca el contenido en menu_versions usando ese current_version_id
		// El contenido obtenido será el menú actual que se pasará al prompt
		menu, err := getMenuById(ctx, input.MenuID, "")
		if err != nil && err != supabaserepo.ErrMenuNotFound {
			obs.Logger.ErrorContext(ctx, "error getting menu from supabase", "error", err, "menuID", input.MenuID)
			return "", err
		}
		
		// Si el menú no se encuentra, continuamos con un menú vacío (comportamiento original)
		if err == supabaserepo.ErrMenuNotFound {
			obs.Logger.WarnContext(ctx, "menu not found in supabase, continuing with empty menu", "menuID", input.MenuID)
		} else {
			obs.Logger.InfoContext(ctx, "menu content retrieved from menu_versions using current_version_id, will be passed to prompt", "menuID", input.MenuID)
		}
		
		// 4. Asignar el contenido del menú obtenido desde menu_versions (usando current_version_id) al request
		// Este contenido será convertido a toon format y pasado al prompt como MENU_ACTUAL
		// Representa la versión actual del menú que el agente debe conocer para aplicar los cambios
		// NOTA: El versionID del contexto se mantiene y se usará cuando se guarde el menú modificado
		input.JsonMenu = menu

		// 5. Llamar al Agente Procesador (que hace la inspección del FunctionCall)
		// El contexto con el versionID se propaga automáticamente y se usará al guardar el menú
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
