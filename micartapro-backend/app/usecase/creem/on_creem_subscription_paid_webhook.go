package creem

import (
	"context"
	"encoding/json"
	"time"

	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"micartapro/app/usecase/billing"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/google/uuid"
)

type OnCreemSubscriptionPaidWebhook func(ctx context.Context, input events.CreemSubscriptionPaidWebhook) error

func init() {
	ioc.Register(NewOnCreemSubscriptionPaidWebhook)
}

func NewOnCreemSubscriptionPaidWebhook(
	obs observability.Observability,
	saveBillingEvent supabaserepo.SaveBillingEvent,
	saveSubscription supabaserepo.SaveSubscription) OnCreemSubscriptionPaidWebhook {
	return func(ctx context.Context, input events.CreemSubscriptionPaidWebhook) error {
		spanCtx, span := obs.Tracer.Start(ctx, "on_creem_subscription_paid_webhook")
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

		// Guardar la subscription si tenemos userID
		if userID != nil {
			subscription := mapCreemSubscriptionPaidWebhookToSubscription(input, *userID)
			obs.Logger.InfoContext(spanCtx, "saving_subscription", "subscriptionID", subscription.SubscriptionID)

			if err := saveSubscription(spanCtx, subscription); err != nil {
				obs.Logger.ErrorContext(spanCtx, "error_saving_subscription", "error", err)
				return err
			}

			obs.Logger.InfoContext(spanCtx, "subscription_saved_successfully", "subscriptionID", subscription.SubscriptionID)
		}

		return nil
	}
}

// mapCreemSubscriptionPaidWebhookToSubscription convierte un webhook de subscription.paid a billing.Subscription
func mapCreemSubscriptionPaidWebhookToSubscription(webhook events.CreemSubscriptionPaidWebhook, userID uuid.UUID) billing.Subscription {
	subscription := billing.Subscription{
		UserID:         userID,
		Provider:       "creem",
		SubscriptionID: webhook.Object.ID,
		CustomerID:     webhook.Object.Customer.ID,
		ProductID:      webhook.Object.Product.ID,
		Status:         webhook.Object.Status,
		Metadata:       make(map[string]any),
	}

	// Parsear fechas
	if webhook.Object.CurrentPeriodStartDate != "" {
		if t, err := time.Parse(time.RFC3339, webhook.Object.CurrentPeriodStartDate); err == nil {
			subscription.CurrentPeriodStart = &t
		}
	}
	if webhook.Object.CurrentPeriodEndDate != "" {
		if t, err := time.Parse(time.RFC3339, webhook.Object.CurrentPeriodEndDate); err == nil {
			subscription.CurrentPeriodEnd = &t
		}
	}
	if webhook.Object.CanceledAt != nil {
		if canceledStr, ok := webhook.Object.CanceledAt.(string); ok && canceledStr != "" {
			if t, err := time.Parse(time.RFC3339, canceledStr); err == nil {
				subscription.CanceledAt = &t
			}
		}
	}

	// Parsear CreatedAt y UpdatedAt
	if webhook.Object.CreatedAt != "" {
		if t, err := time.Parse(time.RFC3339, webhook.Object.CreatedAt); err == nil {
			subscription.CreatedAt = t
		}
	}
	if webhook.Object.UpdatedAt != "" {
		if t, err := time.Parse(time.RFC3339, webhook.Object.UpdatedAt); err == nil {
			subscription.UpdatedAt = t
		}
	}

	// Convertir metadata
	if webhook.Object.Metadata.UserID != "" {
		subscription.Metadata["user_id"] = webhook.Object.Metadata.UserID
	}

	return subscription
}
