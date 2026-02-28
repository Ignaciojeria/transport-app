package supabaserepo

import (
	"context"
	"encoding/json"
	"errors"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/supabase-community/supabase-go"
)

var ErrSlugNotFound = errors.New("slug not found")

type MenuSlugInfo struct {
	UserID string
	MenuID string
}

type GetMenuSlugBySlug func(ctx context.Context, slug string) (MenuSlugInfo, error)

func init() {
	ioc.Register(NewGetMenuSlugBySlug)
}

func NewGetMenuSlugBySlug(supabase *supabase.Client) GetMenuSlugBySlug {
	return func(ctx context.Context, slug string) (MenuSlugInfo, error) {
		// Primero obtener el menu_id desde menu_slugs
		var slugResult []struct {
			MenuID string `json:"menu_id"`
		}

		data, _, err := supabase.From("menu_slugs").
			Select("menu_id", "", false).
			Eq("slug", slug).
			Eq("is_active", "true").
			Execute()

		if err != nil {
			// Verificar si es un error de "no encontrado"
			if err.Error() == "PGRST116" || err.Error() == "no rows in result set" {
				return MenuSlugInfo{}, ErrSlugNotFound
			}
			return MenuSlugInfo{}, err
		}

		if err := json.Unmarshal(data, &slugResult); err != nil {
			return MenuSlugInfo{}, err
		}

		if len(slugResult) == 0 {
			return MenuSlugInfo{}, ErrSlugNotFound
		}

		// Luego obtener el user_id desde user_menus usando el menu_id
		var userMenuResult []struct {
			UserID string `json:"user_id"`
		}

		data, _, err = supabase.From("user_menus").
			Select("user_id", "", false).
			Eq("menu_id", slugResult[0].MenuID).
			Execute()

		if err != nil {
			// Verificar si es un error de "no encontrado"
			if err.Error() == "PGRST116" || err.Error() == "no rows in result set" {
				return MenuSlugInfo{}, ErrSlugNotFound
			}
			return MenuSlugInfo{}, err
		}

		if err := json.Unmarshal(data, &userMenuResult); err != nil {
			return MenuSlugInfo{}, err
		}

		if len(userMenuResult) == 0 {
			return MenuSlugInfo{}, ErrSlugNotFound
		}

		return MenuSlugInfo{
			UserID: userMenuResult[0].UserID,
			MenuID: slugResult[0].MenuID,
		}, nil
	}
}
