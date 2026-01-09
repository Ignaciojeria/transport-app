package supabaserepo

import (
	"context"
	"encoding/json"

	"micartapro/app/shared/infrastructure/supabasecli"
	"micartapro/app/usecase/billing"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/supabase-community/supabase-go"
)

type SaveBillingEvent func(ctx context.Context, billingEvent billing.BillingEvent) error

func init() {
	ioc.Registry(NewSaveBillingEvent, supabasecli.NewSupabaseClient)
}

func NewSaveBillingEvent(supabase *supabase.Client) SaveBillingEvent {
	return func(ctx context.Context, billingEvent billing.BillingEvent) error {
		// Preparar el registro para insertar/actualizar
		record := map[string]interface{}{
			"provider":            billingEvent.Provider,
			"provider_event_id":   billingEvent.ProviderEventID,
			"event_type":          billingEvent.EventType,
			"subscription_id":     billingEvent.SubscriptionID,
			"payload":             json.RawMessage(billingEvent.Payload),
			"provider_created_at": billingEvent.ProviderCreatedAt.Format("2006-01-02T15:04:05Z07:00"), // RFC3339
		}

		// Agregar user_id solo si existe
		if billingEvent.UserID != nil {
			record["user_id"] = billingEvent.UserID.String()
		}

		// Hacer upsert usando (provider, provider_event_id) como clave única para idempotencia
		_, _, err := supabase.From("billing_events").
			Upsert(record, "provider,provider_event_id", "", "").
			Execute()

		return err
	}
}
