package supabaserepo

import (
	"context"
	"encoding/json"
	"fmt"

	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/supabase-community/supabase-go"
)

type UpdateMenuPresentationStyle func(ctx context.Context, menuID string, style events.MenuPresentationStyle) error

func init() {
	ioc.Register(NewUpdateMenuPresentationStyle)
}

func NewUpdateMenuPresentationStyle(supabaseClient *supabase.Client, obs observability.Observability) UpdateMenuPresentationStyle {
	return func(ctx context.Context, menuID string, style events.MenuPresentationStyle) error {
		// 1. Obtener current_version_id del menú
		var menuResult []struct {
			CurrentVersionID *string `json:"current_version_id"`
		}
		data, _, err := supabaseClient.From("menus").
			Select("current_version_id", "", false).
			Eq("id", menuID).
			Execute()
		if err != nil {
			if err.Error() == "PGRST116" || err.Error() == "no rows in result set" {
				return ErrMenuNotFound
			}
			return fmt.Errorf("error querying menus: %w", err)
		}
		if err := json.Unmarshal(data, &menuResult); err != nil || len(menuResult) == 0 {
			return ErrMenuNotFound
		}
		if menuResult[0].CurrentVersionID == nil || *menuResult[0].CurrentVersionID == "" {
			return fmt.Errorf("menu has no current version")
		}
		versionID := *menuResult[0].CurrentVersionID

		// 2. Obtener contenido actual de la versión
		var versionResult []struct {
			Content json.RawMessage `json:"content"`
		}
		data, _, err = supabaseClient.From("menu_versions").
			Select("content", "", false).
			Eq("id", versionID).
			Eq("menu_id", menuID).
			Execute()
		if err != nil || len(data) == 0 {
			if err != nil && (err.Error() == "PGRST116" || err.Error() == "no rows in result set") {
				return ErrMenuNotFound
			}
			return fmt.Errorf("error querying menu_versions: %w", err)
		}
		if err := json.Unmarshal(data, &versionResult); err != nil || len(versionResult) == 0 {
			return ErrMenuNotFound
		}

		var menu events.MenuCreateRequest
		if err := json.Unmarshal(versionResult[0].Content, &menu); err != nil {
			return fmt.Errorf("error unmarshaling menu content: %w", err)
		}
		menu.PresentationStyle = style
		menu.EnsurePresentationStyleDefault()
		contentBytes, err := json.Marshal(menu)
		if err != nil {
			return fmt.Errorf("error marshaling menu content: %w", err)
		}

		// 3. Actualizar content en menu_versions
		_, _, err = supabaseClient.From("menu_versions").
			Update(map[string]interface{}{"content": json.RawMessage(contentBytes)}, "", "").
			Eq("id", versionID).
			Eq("menu_id", menuID).
			Execute()
		if err != nil {
			return fmt.Errorf("error updating menu_versions content: %w", err)
		}
		obs.Logger.InfoContext(ctx, "menu presentation style updated", "menuID", menuID, "versionID", versionID, "style", string(style))
		return nil
	}
}
