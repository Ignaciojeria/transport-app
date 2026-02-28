package fuegoapi

import (
	"encoding/json"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/usecase/mercadopago"
	"net/http"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Register(mercadopagoWebhook)
}

func mercadopagoWebhook(
	s httpserver.Server,
	obs observability.Observability,
	processWebhook mercadopago.ProcessMercadoPagoWebhook,
) {
	fuego.Post(s.Manager, "/webhooks/mercadopago",
		func(c fuego.ContextWithBody[any]) (map[string]interface{}, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "mercadopagoWebhook")
			defer span.End()

			// Obtener el body del request
			requestBodyAny, err := c.Body()
			if err != nil {
				return nil, fuego.HTTPError{
					Title:  "error getting request body",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}

			// Convertir a JSON para procesar
			bodyBytes, err := json.Marshal(requestBodyAny)
			if err != nil {
				return nil, fuego.HTTPError{
					Title:  "error marshaling request body",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}

			// Parsear como map[string]interface{}
			var webhookData map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &webhookData); err != nil {
				return nil, fuego.HTTPError{
					Title:  "error unmarshaling request body",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}

			obs.Logger.InfoContext(spanCtx, "received mercado pago webhook", "webhookData", webhookData)

			// Procesar el webhook
			if err := processWebhook(spanCtx, webhookData); err != nil {
				obs.Logger.ErrorContext(spanCtx, "error processing mercado pago webhook", "error", err)
				return nil, fuego.HTTPError{
					Title:  "error processing webhook",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			obs.Logger.InfoContext(spanCtx, "mercadopago webhook processed successfully")

			// Mercado Pago espera un 200 OK como respuesta
			return map[string]interface{}{
				"status": "ok",
			}, nil
		},
		option.Summary("mercadopagoWebhook"),
		option.Tags("webhooks", "mercadopago"),
	)
}
