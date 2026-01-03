package usecase

import (
	"context"
	"micartapro/app/adapter/out/storage"
	"micartapro/app/domain"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type OnCreemSubscriptionTrialingWebhook func(ctx context.Context, input events.CreemSubscriptionTrialingWebhook) error

func init() {
	ioc.Registry(NewOnCreemSubscriptionTrialingWebhook,
		observability.NewObservability,
		storage.NewSaveEntitlement)
}

func NewOnCreemSubscriptionTrialingWebhook(
	obs observability.Observability,
	saveEntitlement storage.SaveEntitlement) OnCreemSubscriptionTrialingWebhook {
	return func(ctx context.Context, input events.CreemSubscriptionTrialingWebhook) error {
		spanCtx, span := obs.Tracer.Start(ctx, "on_creem_subscription_trialing_webhook")
		defer span.End()

		// El controlador ya validó que el userID existe antes de publicar el evento
		// El subscriber ya lo restauró en el contexto desde el CloudEvent
		userID, _ := sharedcontext.UserIDFromContext(spanCtx)

		// Determinar access basado en el status
		access := false
		switch input.Object.Status {
		case "trialing", "active":
			access = true
		}

		// Mapear el webhook al dominio Entitlement
		entitlement := domain.Entitlement{
			V:        1,
			UserID:   userID,
			Plan:     "pro",
			Status:   input.Object.Status,
			Access:   access,
			StartsAt: input.Object.CurrentPeriodStartDate,
			EndsAt:   input.Object.CurrentPeriodEndDate,
		}

		obs.Logger.InfoContext(spanCtx, "saving_entitlement_from_webhook", "userID", userID)

		if err := saveEntitlement(spanCtx, entitlement); err != nil {
			obs.Logger.ErrorContext(spanCtx, "error_saving_entitlement", "error", err)
			return err
		}

		obs.Logger.InfoContext(spanCtx, "entitlement_saved_successfully", "userID", userID)
		return nil
	}
}
