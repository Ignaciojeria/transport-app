package supabaserepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"micartapro/app/shared/infrastructure/supabasecli"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
)

var ErrCustomerNotFound = errors.New("customer not found")

type GetCustomerID func(ctx context.Context, userID uuid.UUID) (string, error)

func init() {
	ioc.Registry(NewGetCustomerID, supabasecli.NewSupabaseClient)
}

func NewGetCustomerID(supabase *supabase.Client) GetCustomerID {
	return func(ctx context.Context, userID uuid.UUID) (string, error) {
		var result []struct {
			CustomerID string `json:"customer_id"`
		}

		data, _, err := supabase.From("subscriptions").
			Select("customer_id", "", false).
			Eq("user_id", userID.String()).
			Execute()

		if err != nil {
			if err.Error() == "PGRST116" || err.Error() == "no rows in result set" {
				return "", ErrCustomerNotFound
			}
			return "", fmt.Errorf("error querying subscriptions: %w", err)
		}

		if err := json.Unmarshal(data, &result); err != nil {
			return "", fmt.Errorf("error unmarshaling subscription result: %w", err)
		}

		if len(result) == 0 || result[0].CustomerID == "" {
			return "", ErrCustomerNotFound
		}

		return result[0].CustomerID, nil
	}
}
