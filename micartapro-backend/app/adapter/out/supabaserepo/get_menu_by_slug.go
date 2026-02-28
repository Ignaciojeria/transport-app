package supabaserepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"micartapro/app/events"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/supabase-community/supabase-go"
)

var ErrMenuNotFound = errors.New("menu not found")

type GetMenuBySlug func(ctx context.Context, slug string, versionID string) (events.MenuCreateRequest, error)

func init() {
	ioc.Register(NewGetMenuBySlug)
}

func NewGetMenuBySlug(supabase *supabase.Client) GetMenuBySlug {
	return func(ctx context.Context, slug string, versionID string) (events.MenuCreateRequest, error) {
		// 1. Obtener menu_id desde menu_slugs usando el slug
		var slugResult []struct {
			MenuID string `json:"menu_id"`
		}

		data, _, err := supabase.From("menu_slugs").
			Select("menu_id", "", false).
			Eq("slug", slug).
			Eq("is_active", "true").
			Execute()

		if err != nil {
			if err.Error() == "PGRST116" || err.Error() == "no rows in result set" {
				return events.MenuCreateRequest{}, ErrMenuNotFound
			}
			return events.MenuCreateRequest{}, fmt.Errorf("error querying menu_slugs: %w", err)
		}

		if err := json.Unmarshal(data, &slugResult); err != nil {
			return events.MenuCreateRequest{}, fmt.Errorf("error unmarshaling slug result: %w", err)
		}

		if len(slugResult) == 0 {
			return events.MenuCreateRequest{}, ErrMenuNotFound
		}

		menuID := slugResult[0].MenuID

		// 2. Obtener el menú desde menus (incluye user_id y current_version_id)
		var menuResult []struct {
			UserID           string  `json:"user_id"`
			CurrentVersionID *string `json:"current_version_id"`
		}

		data, _, err = supabase.From("menus").
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

		// Determinar qué versión usar: si se proporciona version_id, usarlo; si no, usar current_version_id
		var targetVersionID string
		if versionID != "" {
			// Validar que el version_id proporcionado pertenezca al menú
			targetVersionID = versionID
		} else {
			// Usar la versión actual si no se especifica una versión
			if menuResult[0].CurrentVersionID == nil || *menuResult[0].CurrentVersionID == "" {
				return events.MenuCreateRequest{}, fmt.Errorf("menu has no current version")
			}
			targetVersionID = *menuResult[0].CurrentVersionID
		}

		// 3. Obtener el contenido de la versión especificada (o actual) desde menu_versions
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
				return events.MenuCreateRequest{}, ErrMenuNotFound
			}
			return events.MenuCreateRequest{}, fmt.Errorf("error querying menu_versions: %w", err)
		}

		if err := json.Unmarshal(data, &versionResult); err != nil {
			return events.MenuCreateRequest{}, fmt.Errorf("error unmarshaling version result: %w", err)
		}

		if len(versionResult) == 0 {
			return events.MenuCreateRequest{}, ErrMenuNotFound
		}

		// 4. Deserializar el contenido del menú
		var menu events.MenuCreateRequest
		if err := json.Unmarshal(versionResult[0].Content, &menu); err != nil {
			return events.MenuCreateRequest{}, fmt.Errorf("error unmarshaling menu content: %w", err)
		}

		// Asegurar que el ID del menú esté establecido
		menu.ID = menuID
		menu.EnsurePresentationStyleDefault()

		return menu, nil
	}
}
