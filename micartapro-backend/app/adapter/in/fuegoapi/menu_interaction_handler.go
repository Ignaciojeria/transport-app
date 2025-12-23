package fuegoapi

import (
	"micartapro/app/domain"
	"micartapro/app/shared/infrastructure/auth"
	"micartapro/app/shared/infrastructure/eventprocessing"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"net/http"
	"strings"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(
		menuInteractionHandler,
		httpserver.New,
		observability.NewObservability,
		eventprocessing.NewPublisherStrategy,
		auth.NewSupabaseTokenValidator,
	)
}
func menuInteractionHandler(
	s httpserver.Server,
	obs observability.Observability,
	publisherManager eventprocessing.PublisherManager,
	supabaseTokenValidator auth.SupabaseTokenValidator,
) {
	fuego.Post(s.Manager, "/menu/interaction",
		func(c fuego.ContextWithBody[domain.MenuInteractionRequest]) (any, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "menuInteractionHandler")
			defer span.End()

			token := c.Header("Authorization")
			token = strings.TrimPrefix(token, "Bearer ")
			_, err := supabaseTokenValidator.ValidateJWT(token)
			if err != nil {
				return nil, fuego.HTTPError{
					Title:  "error validating token",
					Detail: err.Error(),
					Status: http.StatusUnauthorized,
				}
			}
			body, err := c.Body()
			if err != nil {
				return nil, fuego.HTTPError{
					Title:  "error getting request body",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			if err := body.Validate(); err != nil {
				return nil, fuego.HTTPError{
					Title:  "error validating request body",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}
			if err := publisherManager.Publish(spanCtx, eventprocessing.PublishRequest{
				Topic:  "micartapro.events",
				Source: "micartapro.api.menu.interaction",
				Event:  body,
			}); err != nil {
				return nil, fuego.HTTPError{
					Title:  "error publishing event",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			obs.Logger.InfoContext(spanCtx, "menuInteractionRequest published", "requestBody", body)
			return http.StatusOK, nil
		}, option.Summary("menuInteractionRequest"), option.Header("Idempotency-Key", "01KCW67YKSV455GBVDT88S4072"))
}
