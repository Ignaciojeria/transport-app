package supabaserepo

import (
	"context"
	"encoding/json"
	"fmt"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/supabase-community/supabase-go"
)

// GetUserMenuIds retorna todos los menu_id del usuario (para validar acceso).
type GetUserMenuIds func(ctx context.Context, userID string) ([]string, error)

func init() {
	ioc.Register(NewGetUserMenuIds)
}

func NewGetUserMenuIds(sb *supabase.Client) GetUserMenuIds {
	return func(ctx context.Context, userID string) ([]string, error) {
		var result []struct {
			MenuID string `json:"menu_id"`
		}

		data, _, err := sb.From("user_menus").
			Select("menu_id", "", false).
			Eq("user_id", userID).
			Execute()

		if err != nil {
			return nil, fmt.Errorf("error querying user_menus: %w", err)
		}

		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("error unmarshaling user_menus result: %w", err)
		}

		ids := make([]string, 0, len(result))
		for _, r := range result {
			if r.MenuID != "" {
				ids = append(ids, r.MenuID)
			}
		}
		return ids, nil
	}
}
