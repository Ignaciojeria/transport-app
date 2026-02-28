package creem

import (
	"context"
	"encoding/json"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/google/uuid"
)

type OnCreemSubscriptionUpdateWebhook func(ctx context.Context, input events.CreemSubscriptionUpdateWebhook) error

func init() {
	ioc.Register(NewOnCreemSubscriptionUpdateWebhook)
}

func NewOnCreemSubscriptionUpdateWebhook(
	obs observability.Observability,
	saveBillingEvent supabaserepo.SaveBillingEvent) OnCreemSubscriptionUpdateWebhook {
	return func(ctx context.Context, input events.CreemSubscriptionUpdateWebhook) error {
		spanCtx, span := obs.Tracer.Start(ctx, "on_creem_subscription_update_webhook")
		defer span.End()

		userIDStr, _ := sharedcontext.UserIDFromContext(spanCtx)
		var userID *uuid.UUID
		if userIDStr != "" {
			if parsed, err := uuid.Parse(userIDStr); err == nil {
				userID = &parsed
			}
		}

		payload, err := json.Marshal(input)
		if err != nil {
			obs.Logger.ErrorContext(spanCtx, "error_marshaling_webhook", "error", err)
			return err
		}

		billingEvent := toBillingEvent(
			input.ID,
			input.EventType,
			input.CreatedAt,
			json.RawMessage(payload),
			userID,
		)

		obs.Logger.InfoContext(spanCtx, "saving_billing_event", "eventID", input.ID, "eventType", input.EventType)

		if err := saveBillingEvent(spanCtx, billingEvent); err != nil {
			obs.Logger.ErrorContext(spanCtx, "error_saving_billing_event", "error", err)
			return err
		}

		obs.Logger.InfoContext(spanCtx, "billing_event_saved_successfully", "eventID", input.ID)
		return nil
	}
}
