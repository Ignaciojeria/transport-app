package supabaserepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"micartapro/app/usecase/billing"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
)

func init() {
	ioc.Register(NewGetUserCredits)
}

func NewGetUserCredits(supabase *supabase.Client) GetUserCredits {
	return func(ctx context.Context, userID uuid.UUID) (*billing.UserCredits, error) {
		var result []struct {
			ID        int64     `json:"id"`
			UserID    string    `json:"user_id"`
			Balance   int       `json:"balance"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		}

		data, _, err := supabase.From("user_credits").
			Select("*", "", false).
			Eq("user_id", userID.String()).
			Execute()

		if err != nil {
			return nil, fmt.Errorf("error querying user_credits: %w", err)
		}

		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("error unmarshaling user_credits result: %w", err)
		}

		if len(result) == 0 {
			// Si no existe, retornar con balance 0
			return &billing.UserCredits{
				UserID:  userID,
				Balance: 0,
			}, nil
		}

		parsedUserID, err := uuid.Parse(result[0].UserID)
		if err != nil {
			return nil, fmt.Errorf("error parsing user_id: %w", err)
		}

		return &billing.UserCredits{
			UserID:    parsedUserID,
			Balance:   result[0].Balance,
			CreatedAt: result[0].CreatedAt,
			UpdatedAt: result[0].UpdatedAt,
		}, nil
	}
}
