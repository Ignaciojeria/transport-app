package menu

import (
	"context"
	"fmt"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/eventprocessing"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"micartapro/app/usecase/billing"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/google/uuid"
)

type OnUserMenusInsertedWebhook func(ctx context.Context, input events.UserMenusInsertedWebhook) error

func init() {
	ioc.Registry(NewOnUserMenusInsertedWebhook, eventprocessing.NewPublisherStrategy, supabaserepo.NewGetUserCredits, supabaserepo.NewGrantCredits, observability.NewObservability)
}

func NewOnUserMenusInsertedWebhook(
	publisherManager eventprocessing.PublisherManager,
	getUserCredits supabaserepo.GetUserCredits,
	grantCredits supabaserepo.GrantCredits,
	obs observability.Observability) OnUserMenusInsertedWebhook {
	return func(ctx context.Context, wh events.UserMenusInsertedWebhook) error {
		spanCtx, span := obs.Tracer.Start(ctx, "on_user_menus_inserted_webhook")
		defer span.End()

		// Parsear el user_id del webhook
		userID, err := uuid.Parse(wh.Record.UserID)
		if err != nil {
			obs.Logger.ErrorContext(spanCtx, "error_parsing_user_id", "error", err, "user_id", wh.Record.UserID)
			return fmt.Errorf("error parsing user_id: %w", err)
		}

		// Verificar si el usuario ya tiene créditos (para evitar otorgar créditos duplicados)
		userCredits, errCredits := getUserCredits(spanCtx, userID)
		if errCredits != nil {
			obs.Logger.ErrorContext(spanCtx, "error_getting_user_credits", "error", errCredits, "user_id", userID)
			// Continuar aunque haya error, ya que podría ser que no exista registro aún
		}

		// Otorgar 2 créditos de prueba solo si el usuario no tiene créditos o tiene 0
		if userCredits == nil || userCredits.Balance == 0 {
			menuIDStr := wh.Record.MenuID
			description := "Welcome bonus: 2 trial credits for new user registration"
			_, err := grantCredits(spanCtx, billing.GrantCreditsRequest{
				UserID:      userID,
				Amount:      2,
				Source:      "registration.welcome_bonus",
				SourceID:    &menuIDStr, // Usar menu_id como source_id para evitar duplicados
				Description: &description,
			})
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error_granting_welcome_credits", "error", err, "user_id", userID)
				// No retornar error para no bloquear el flujo principal, solo loguear
			} else {
				obs.Logger.InfoContext(spanCtx, "welcome_credits_granted", "user_id", userID, "amount", 2)
			}
		} else {
			obs.Logger.InfoContext(spanCtx, "user_already_has_credits", "user_id", userID, "balance", userCredits.Balance)
		}

		// Generar version_id si no viene en el contexto
		var versionID string
		if existingVersionID, ok := sharedcontext.VersionIDFromContext(spanCtx); ok && existingVersionID != "" {
			versionID = existingVersionID
		} else {
			versionID = uuid.New().String()
			spanCtx = sharedcontext.WithVersionID(spanCtx, versionID)
		}

		menuCreateRequest := events.MenuCreateRequest{
			ID:                wh.Record.MenuID,
			PresentationStyle: events.MenuStyleHero,
			CoverImage:        "https://storage.googleapis.com/micartapro-menus/core/micartaprov3.webp",
			BusinessInfo: events.BusinessInfo{
				BusinessName:  "cadorago",
				Whatsapp:      "+56957857558",
				BusinessHours: []string{},
				Currency:      events.CurrencyCLP,
			},
			Menu: []events.MenuCategory{},
			DeliveryOptions: []events.DeliveryOption{
				{
					Type:        events.DeliveryOptionDelivery,
					RequireTime: false,
				},
				{
					Type:            events.DeliveryOptionPickup,
					RequireTime:     true,
					TimeRequestType: events.TimeRequestWindow,
					TimeWindows: []events.TimeWindow{
						{
							Start: "09:00",
							End:   "23:59",
						},
					},
				},
			},
		}

		err = publisherManager.Publish(spanCtx, eventprocessing.PublishRequest{
			Topic:       "micartapro.events",
			Source:      "micartapro.webhook.user.menus.inserted",
			OrderingKey: wh.Record.MenuID,
			Event:       menuCreateRequest,
		})
		if err != nil {
			return err
		}
		return nil
	}
}
