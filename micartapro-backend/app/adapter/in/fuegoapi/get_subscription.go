package fuegoapi

import (
	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/google/uuid"
)

func init() {
	ioc.Registry(
		getSubscription,
		httpserver.New,
		supabaserepo.NewGetSubscription,
		observability.NewObservability,
		apimiddleware.NewJWTAuthMiddleware)
}

type SubscriptionResponse struct {
	UserID             string                 `json:"user_id"`
	Provider           string                 `json:"provider"`
	SubscriptionID     string                 `json:"subscription_id"`
	CustomerID         string                 `json:"customer_id"`
	ProductID          string                 `json:"product_id"`
	Status             string                 `json:"status"`
	CurrentPeriodStart *string                `json:"current_period_start,omitempty"`
	CurrentPeriodEnd   *string                `json:"current_period_end,omitempty"`
	CancelAt           *string                `json:"cancel_at,omitempty"`
	CanceledAt         *string                `json:"canceled_at,omitempty"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt          string                 `json:"created_at"`
	UpdatedAt          string                 `json:"updated_at"`
	HasActiveSubscription bool                `json:"has_active_subscription"`
}

func getSubscription(
	s httpserver.Server,
	getSubscription supabaserepo.GetSubscription,
	obs observability.Observability,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware) {
	fuego.Get(s.Manager, "/subscription",
		func(c fuego.ContextNoBody) (any, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "getSubscription")
			defer span.End()

			// Extraer user_id del contexto
			userIDStr, ok := sharedcontext.UserIDFromContext(spanCtx)
			if !ok || userIDStr == "" {
				return nil, fuego.UnauthorizedError{
					Title:  "user_id not found in context",
					Detail: "user_id not found in context",
					Status: 401,
				}
			}

			// Parsear userID
			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				return nil, fuego.BadRequestError{
					Title:  "invalid user_id",
					Detail: "invalid user_id format",
					Status: 400,
				}
			}

			// Obtener la suscripción
			subscription, err := getSubscription(spanCtx, userID)
			if err != nil {
				if err == supabaserepo.ErrSubscriptionNotFound {
					// No hay suscripción, devolver respuesta indicando que no tiene suscripción activa
					return SubscriptionResponse{
						HasActiveSubscription: false,
					}, nil
				}
				obs.Logger.ErrorContext(spanCtx, "error_getting_subscription", "error", err, "user_id", userIDStr)
				return nil, fuego.HTTPError{
					Title:  "error getting subscription",
					Detail: err.Error(),
					Status: 500,
				}
			}

			// Construir la respuesta
			response := SubscriptionResponse{
				UserID:             subscription.UserID.String(),
				Provider:           subscription.Provider,
				SubscriptionID:     subscription.SubscriptionID,
				CustomerID:         subscription.CustomerID,
				ProductID:          subscription.ProductID,
				Status:             subscription.Status,
				Metadata:           subscription.Metadata,
				CreatedAt:          subscription.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt:          subscription.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
				HasActiveSubscription: subscription.Status == "active" || subscription.Status == "trialing",
			}

			// Agregar fechas opcionales si existen
			if subscription.CurrentPeriodStart != nil {
				startStr := subscription.CurrentPeriodStart.Format("2006-01-02T15:04:05Z07:00")
				response.CurrentPeriodStart = &startStr
			}

			if subscription.CurrentPeriodEnd != nil {
				endStr := subscription.CurrentPeriodEnd.Format("2006-01-02T15:04:05Z07:00")
				response.CurrentPeriodEnd = &endStr
			}

			if subscription.CancelAt != nil {
				cancelAtStr := subscription.CancelAt.Format("2006-01-02T15:04:05Z07:00")
				response.CancelAt = &cancelAtStr
			}

			if subscription.CanceledAt != nil {
				canceledAtStr := subscription.CanceledAt.Format("2006-01-02T15:04:05Z07:00")
				response.CanceledAt = &canceledAtStr
			}

			return response, nil
		}, option.Summary("getSubscription"), option.Middleware(jwtAuthMiddleware))
}
