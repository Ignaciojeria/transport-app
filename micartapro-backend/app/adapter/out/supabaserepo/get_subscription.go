package supabaserepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"micartapro/app/usecase/billing"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
)

var ErrSubscriptionNotFound = errors.New("subscription not found")

type GetSubscription func(ctx context.Context, userID uuid.UUID) (*billing.Subscription, error)

func init() {
	ioc.Register(NewGetSubscription)
}

func NewGetSubscription(supabase *supabase.Client) GetSubscription {
	return func(ctx context.Context, userID uuid.UUID) (*billing.Subscription, error) {
		var result []struct {
			ID                 int64           `json:"id"`
			UserID             string          `json:"user_id"`
			Provider           string          `json:"provider"`
			SubscriptionID     string          `json:"subscription_id"`
			CustomerID         string          `json:"customer_id"`
			ProductID          string          `json:"product_id"`
			Status             string          `json:"status"`
			CurrentPeriodStart *string         `json:"current_period_start"`
			CurrentPeriodEnd   *string         `json:"current_period_end"`
			CancelAt           *string         `json:"cancel_at"`
			CanceledAt         *string         `json:"canceled_at"`
			Metadata           json.RawMessage `json:"metadata"`
			CreatedAt          time.Time       `json:"created_at"`
			UpdatedAt          time.Time       `json:"updated_at"`
		}

		data, _, err := supabase.From("subscriptions").
			Select("*", "", false).
			Eq("user_id", userID.String()).
			Execute()

		if err != nil {
			if err.Error() == "PGRST116" || err.Error() == "no rows in result set" {
				return nil, ErrSubscriptionNotFound
			}
			return nil, fmt.Errorf("error querying subscriptions: %w", err)
		}

		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("error unmarshaling subscription result: %w", err)
		}

		if len(result) == 0 {
			return nil, ErrSubscriptionNotFound
		}

		sub := result[0]

		// Parsear userID
		parsedUserID, err := uuid.Parse(sub.UserID)
		if err != nil {
			return nil, fmt.Errorf("error parsing user_id: %w", err)
		}

		// Construir la suscripción
		subscription := &billing.Subscription{
			UserID:         parsedUserID,
			Provider:       sub.Provider,
			SubscriptionID: sub.SubscriptionID,
			CustomerID:     sub.CustomerID,
			ProductID:      sub.ProductID,
			Status:         sub.Status,
			CreatedAt:      sub.CreatedAt,
			UpdatedAt:      sub.UpdatedAt,
		}

		// Parsear fechas opcionales
		if sub.CurrentPeriodStart != nil {
			t, err := time.Parse(time.RFC3339, *sub.CurrentPeriodStart)
			if err == nil {
				subscription.CurrentPeriodStart = &t
			}
		}

		if sub.CurrentPeriodEnd != nil {
			t, err := time.Parse(time.RFC3339, *sub.CurrentPeriodEnd)
			if err == nil {
				subscription.CurrentPeriodEnd = &t
			}
		}

		if sub.CancelAt != nil {
			t, err := time.Parse(time.RFC3339, *sub.CancelAt)
			if err == nil {
				subscription.CancelAt = &t
			}
		}

		if sub.CanceledAt != nil {
			t, err := time.Parse(time.RFC3339, *sub.CanceledAt)
			if err == nil {
				subscription.CanceledAt = &t
			}
		}

		// Parsear metadata
		if len(sub.Metadata) > 0 {
			var metadata map[string]interface{}
			if err := json.Unmarshal(sub.Metadata, &metadata); err == nil {
				subscription.Metadata = metadata
			} else {
				subscription.Metadata = make(map[string]interface{})
			}
		} else {
			subscription.Metadata = make(map[string]interface{})
		}

		return subscription, nil
	}
}
