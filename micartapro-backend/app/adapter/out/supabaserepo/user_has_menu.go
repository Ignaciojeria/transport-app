package supabaserepo

import (
	"context"
	"encoding/json"
	"fmt"

	"micartapro/app/shared/infrastructure/supabasecli"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/supabase-community/supabase-go"
)

// UserHasMenu verifica si el usuario tiene acceso al menú (existe en user_menus).
type UserHasMenu func(ctx context.Context, userID, menuID string) (bool, error)

func init() {
	ioc.Registry(NewUserHasMenu, supabasecli.NewSupabaseClient)
}

func NewUserHasMenu(supabase *supabase.Client) UserHasMenu {
	return func(ctx context.Context, userID, menuID string) (bool, error) {
		var result []struct {
			MenuID string `json:"menu_id"`
		}

		data, _, err := supabase.From("user_menus").
			Select("menu_id", "", false).
			Eq("user_id", userID).
			Eq("menu_id", menuID).
			Execute()

		if err != nil {
			return false, fmt.Errorf("error querying user_menus: %w", err)
		}

		if err := json.Unmarshal(data, &result); err != nil {
			return false, fmt.Errorf("error unmarshaling user_menus result: %w", err)
		}

		return len(result) > 0, nil
	}
}
