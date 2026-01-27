package supabaserepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"micartapro/app/shared/infrastructure/supabasecli"
	"micartapro/app/usecase/billing"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
)

func init() {
	ioc.Registry(NewGetCreditTransactions, supabasecli.NewSupabaseClient)
}

func NewGetCreditTransactions(supabase *supabase.Client) GetCreditTransactions {
	return func(ctx context.Context, userID uuid.UUID, limit int) ([]billing.CreditTransaction, error) {
		var result []struct {
			ID             int64     `json:"id"`
			UserID         string    `json:"user_id"`
			Amount         int       `json:"amount"`
			TransactionType string   `json:"transaction_type"`
			Source         string    `json:"source"`
			SourceID       *string   `json:"source_id"`
			Description    *string   `json:"description"`
			BalanceBefore  int       `json:"balance_before"`
			BalanceAfter   int       `json:"balance_after"`
			CreatedAt      time.Time `json:"created_at"`
		}

		query := supabase.From("credit_transactions").
			Select("*", "", false).
			Eq("user_id", userID.String())

		if limit > 0 {
			query = query.Limit(limit, "")
		}

		data, _, err := query.Execute()

		if err != nil {
			return nil, fmt.Errorf("error querying credit_transactions: %w", err)
		}

		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("error unmarshaling credit_transactions result: %w", err)
		}

		transactions := make([]billing.CreditTransaction, len(result))
		for i, r := range result {
			parsedUserID, err := uuid.Parse(r.UserID)
			if err != nil {
				return nil, fmt.Errorf("error parsing user_id: %w", err)
			}

			transactions[i] = billing.CreditTransaction{
				ID:              r.ID,
				UserID:          parsedUserID,
				Amount:          r.Amount,
				TransactionType: r.TransactionType,
				Source:          r.Source,
				SourceID:        r.SourceID,
				Description:     r.Description,
				BalanceBefore:   r.BalanceBefore,
				BalanceAfter:    r.BalanceAfter,
				CreatedAt:       r.CreatedAt,
			}
		}

		return transactions, nil
	}
}
