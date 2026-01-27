package supabaserepo

import (
	"context"
	"encoding/json"
	"fmt"
	"micartapro/app/shared/infrastructure/supabasecli"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
)

type HasActiveSubscription func(ctx context.Context, userID uuid.UUID) (bool, error)

func init() {
	ioc.Registry(NewHasActiveSubscription, supabasecli.NewSupabaseClient)
}

func NewHasActiveSubscription(supabase *supabase.Client) HasActiveSubscription {
	return func(ctx context.Context, userID uuid.UUID) (bool, error) {
		var result []struct {
			ID int64 `json:"id"`
		}

		data, _, err := supabase.From("subscriptions").
			Select("id", "", false).
			Eq("user_id", userID.String()).
			Execute()

		if err != nil {
			if err.Error() == "PGRST116" || err.Error() == "no rows in result set" {
				return false, nil // No hay suscripción, pero no es un error
			}
			return false, fmt.Errorf("error querying subscriptions: %w", err)
		}

		if err := json.Unmarshal(data, &result); err != nil {
			return false, fmt.Errorf("error unmarshaling subscription result: %w", err)
		}

		if len(result) == 0 {
			return false, nil // No hay suscripción
		}

		return true, nil // Hay suscripción
	}
}
