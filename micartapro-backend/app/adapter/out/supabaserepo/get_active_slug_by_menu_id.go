package supabaserepo

import (
	"context"
	"encoding/json"
	"fmt"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/supabase-community/supabase-go"
)

// GetActiveSlugByMenuID devuelve el slug activo para un menu_id (backend, sin RLS de usuario).
type GetActiveSlugByMenuID func(ctx context.Context, menuID string) (string, error)

func init() {
	ioc.Register(NewGetActiveSlugByMenuID)
}

func NewGetActiveSlugByMenuID(supabase *supabase.Client) GetActiveSlugByMenuID {
	return func(ctx context.Context, menuID string) (string, error) {
		var result []struct {
			Slug string `json:"slug"`
		}
		data, _, err := supabase.From("menu_slugs").
			Select("slug", "", false).
			Eq("menu_id", menuID).
			Eq("is_active", "true").
			Execute()
		if err != nil {
			if err.Error() == "PGRST116" || err.Error() == "no rows in result set" {
				return "", nil
			}
			return "", fmt.Errorf("querying menu_slugs: %w", err)
		}
		if err := json.Unmarshal(data, &result); err != nil {
			return "", fmt.Errorf("unmarshaling menu_slugs: %w", err)
		}
		if len(result) == 0 {
			return "", nil
		}
		s := result[0].Slug
		if s == "" {
			return "", nil
		}
		return s, nil
	}
}
