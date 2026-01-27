package fuegoapi

import (
	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"micartapro/app/usecase/mercadopago"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/google/uuid"
)

func init() {
	ioc.Registry(
		mercadopagoCheckout,
		httpserver.New,
		observability.NewObservability,
		mercadopago.NewCreateMercadoPagoCheckout,
		apimiddleware.NewJWTAuthMiddleware,
	)
}

func mercadopagoCheckout(
	s httpserver.Server,
	obs observability.Observability,
	createCheckout mercadopago.CreateMercadoPagoCheckout,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware,
) {
	fuego.Get(s.Manager, "/checkout/mercadopago",
		func(c fuego.ContextNoBody) (mercadopago.CreateMercadoPagoCheckoutResult, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "mercadopagoCheckout")
			defer span.End()

			// Extraer user_id del contexto del token JWT (opcional, para metadata)
			userID, _ := sharedcontext.UserIDFromContext(spanCtx)

			// Generar un external_reference único para esta orden
			externalReference := uuid.New().String()
			if userID != "" {
				externalReference = userID + "_" + externalReference
			}

			// Payload hardcodeado
			checkoutRequest := events.MercadoPagoCheckoutRequest{
				Items: []events.MercadoPagoCheckoutItem{
					{
						ProductName: "Suscripción",
						Quantity:    1,
						Unit:        "EACH",
						UnitPrice:   3500,
					},
				},
				Totals: struct {
					Subtotal    int    `json:"subtotal"`
					DeliveryFee int    `json:"deliveryFee"`
					Total       int    `json:"total"`
					Currency    string `json:"currency"`
				}{
					Subtotal:    3500,
					DeliveryFee: 0,
					Total:       3500,
					Currency:    "CLP",
				},
				BusinessInfo: struct {
					BusinessName string `json:"businessName"`
					Whatsapp     string `json:"whatsapp"`
				}{
					BusinessName: "MiCartaPro",
					Whatsapp:     "",
				},
				Fulfillment: struct {
					Type string `json:"type"`
				}{
					Type: "DIGITAL",
				},
			}

			// Llamar al caso de uso
			result, err := createCheckout(spanCtx, checkoutRequest, externalReference)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error creating mercado pago checkout", "error", err, "userID", userID)
				return mercadopago.CreateMercadoPagoCheckoutResult{}, fuego.HTTPError{
					Title:  "error creating mercado pago checkout",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			obs.Logger.InfoContext(spanCtx, "mercadopago checkout created successfully", 
				"userID", userID,
				"preferenceID", result.PreferenceID,
				"checkoutURL", result.CheckoutURL)
			
			return result, nil
		},
		option.Summary("mercadopagoCheckout"),
		option.Tags("payments", "mercadopago"),
		option.Middleware(jwtAuthMiddleware),
	)
}
