package supabaserepo

import (
	"context"
	"encoding/json"

	"micartapro/app/usecase/billing"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/supabase-community/supabase-go"
)

type SaveSubscription func(ctx context.Context, subscription billing.Subscription) error

func init() {
	ioc.Register(NewSaveSubscription)
}

func NewSaveSubscription(supabase *supabase.Client) SaveSubscription {
	return func(ctx context.Context, subscription billing.Subscription) error {
		// Preparar el registro para insertar/actualizar
		record := map[string]interface{}{
			"user_id":         subscription.UserID.String(),
			"provider":        subscription.Provider,
			"subscription_id": subscription.SubscriptionID,
			"customer_id":     subscription.CustomerID,
			"product_id":      subscription.ProductID,
			"status":          subscription.Status,
			"updated_at":      subscription.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"), // RFC3339
		}

		// Agregar fechas solo si existen
		if subscription.CurrentPeriodStart != nil {
			record["current_period_start"] = subscription.CurrentPeriodStart.Format("2006-01-02T15:04:05Z07:00")
		}
		if subscription.CurrentPeriodEnd != nil {
			record["current_period_end"] = subscription.CurrentPeriodEnd.Format("2006-01-02T15:04:05Z07:00")
		}
		if subscription.CancelAt != nil {
			record["cancel_at"] = subscription.CancelAt.Format("2006-01-02T15:04:05Z07:00")
		}
		if subscription.CanceledAt != nil {
			record["canceled_at"] = subscription.CanceledAt.Format("2006-01-02T15:04:05Z07:00")
		}

		// Convertir metadata a JSON
		if subscription.Metadata != nil {
			metadataBytes, err := json.Marshal(subscription.Metadata)
			if err == nil {
				record["metadata"] = json.RawMessage(metadataBytes)
			}
		} else {
			record["metadata"] = json.RawMessage("{}")
		}

		// Agregar created_at solo si está definido (si no, la BD usará el default)
		if !subscription.CreatedAt.IsZero() {
			record["created_at"] = subscription.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		// Hacer upsert usando user_id como clave única (constraint one_subscription_per_user)
		// Esto actualizará la subscription existente si ya hay una para este usuario
		_, _, err := supabase.From("subscriptions").
			Upsert(record, "user_id", "", "").
			Execute()

		return err
	}
}
