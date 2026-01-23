package supabaserepo

import (
	"context"
	"encoding/json"
	"fmt"

	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/infrastructure/supabasecli"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/supabase-community/supabase-go"
)

type GetMenuById func(ctx context.Context, menuID string, versionID string) (events.MenuCreateRequest, error)

func init() {
	ioc.Registry(NewGetMenuById, supabasecli.NewSupabaseClient, observability.NewObservability)
}

func NewGetMenuById(supabase *supabase.Client, obs observability.Observability) GetMenuById {
	return func(ctx context.Context, menuID string, versionID string) (events.MenuCreateRequest, error) {
		// 1. Obtener el menú desde menus para obtener el current_version_id
		// Esto nos permite saber qué versión del menú debemos usar
		var menuResult []struct {
			UserID           string  `json:"user_id"`
			CurrentVersionID *string `json:"current_version_id"`
		}

		data, _, err := supabase.From("menus").
			Select("user_id,current_version_id", "", false).
			Eq("id", menuID).
			Execute()

		if err != nil {
			if err.Error() == "PGRST116" || err.Error() == "no rows in result set" {
				return events.MenuCreateRequest{}, ErrMenuNotFound
			}
			return events.MenuCreateRequest{}, fmt.Errorf("error querying menus: %w", err)
		}

		if err := json.Unmarshal(data, &menuResult); err != nil {
			return events.MenuCreateRequest{}, fmt.Errorf("error unmarshaling menu result: %w", err)
		}

		if len(menuResult) == 0 {
			return events.MenuCreateRequest{}, ErrMenuNotFound
		}

		// 2. Determinar qué versión usar: si se proporciona version_id, usarlo; si no, usar current_version_id
		var targetVersionID string
		if versionID != "" {
			// Validar que el version_id proporcionado pertenezca al menú
			targetVersionID = versionID
			obs.Logger.InfoContext(ctx, "using provided version_id to get content", "menuID", menuID, "versionID", targetVersionID)
		} else {
			// Usar la versión actual si no se especifica una versión
			if menuResult[0].CurrentVersionID == nil || *menuResult[0].CurrentVersionID == "" {
				return events.MenuCreateRequest{}, fmt.Errorf("menu has no current version")
			}
			targetVersionID = *menuResult[0].CurrentVersionID
			obs.Logger.InfoContext(ctx, "using current_version_id from menu to get content", "menuID", menuID, "currentVersionID", targetVersionID)
		}

		// 3. Obtener el contenido del menú desde menu_versions usando el version_id determinado
		// Este es el contenido que se pasará al prompt para que el agente conozca el menú anterior
		var versionResult []struct {
			Content json.RawMessage `json:"content"`
		}

		data, _, err = supabase.From("menu_versions").
			Select("content", "", false).
			Eq("id", targetVersionID).
			Eq("menu_id", menuID).
			Execute()

		if err != nil {
			if err.Error() == "PGRST116" || err.Error() == "no rows in result set" {
				obs.Logger.WarnContext(ctx, "menu version not found", "menuID", menuID, "versionID", targetVersionID)
				return events.MenuCreateRequest{}, ErrMenuNotFound
			}
			return events.MenuCreateRequest{}, fmt.Errorf("error querying menu_versions: %w", err)
		}

		if err := json.Unmarshal(data, &versionResult); err != nil {
			return events.MenuCreateRequest{}, fmt.Errorf("error unmarshaling version result: %w", err)
		}

		if len(versionResult) == 0 {
			obs.Logger.WarnContext(ctx, "menu version content not found", "menuID", menuID, "versionID", targetVersionID)
			return events.MenuCreateRequest{}, ErrMenuNotFound
		}

		// 4. Deserializar el contenido del menú desde menu_versions
		// Este contenido será el que se pase al contexto del prompt
		var menu events.MenuCreateRequest
		if err := json.Unmarshal(versionResult[0].Content, &menu); err != nil {
			return events.MenuCreateRequest{}, fmt.Errorf("error unmarshaling menu content from menu_versions: %w", err)
		}

		// Asegurar que el ID del menú esté establecido
		menu.ID = menuID

		obs.Logger.InfoContext(ctx, "menu content retrieved from menu_versions", "menuID", menuID, "versionID", targetVersionID)
		return menu, nil
	}
}
