package supabaserepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/supabase-community/supabase-go"
)

var ErrUserMenuNotFound = errors.New("user menu not found")

type GetMenuIdByUserId func(ctx context.Context, userID string) (string, error)

func init() {
	ioc.Register(NewGetMenuIdByUserId)
}

func NewGetMenuIdByUserId(supabase *supabase.Client) GetMenuIdByUserId {
	return func(ctx context.Context, userID string) (string, error) {
		var result []struct {
			MenuID string `json:"menu_id"`
		}

		data, _, err := supabase.From("user_menus").
			Select("menu_id", "", false).
			Eq("user_id", userID).
			Execute()

		if err != nil {
			if err.Error() == "PGRST116" || err.Error() == "no rows in result set" {
				return "", ErrUserMenuNotFound
			}
			return "", fmt.Errorf("error querying user_menus: %w", err)
		}

		if err := json.Unmarshal(data, &result); err != nil {
			return "", fmt.Errorf("error unmarshaling user_menus result: %w", err)
		}

		if len(result) == 0 {
			return "", ErrUserMenuNotFound
		}

		return result[0].MenuID, nil
	}
}
