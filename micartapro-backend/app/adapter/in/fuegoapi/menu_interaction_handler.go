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
	supabaseTokenValidator auth.SupabaseTokenValidator) {
	fuego.Post(s.Manager, "/menu/interaction",
		func(c fuego.ContextWithBody[domain.MenuInteractionRequest]) (any, error) {
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
			spanCtx, span := obs.Tracer.Start(c.Context(), "menuInteractionRequest")
			defer span.End()
			body, err := c.Body()
			if err != nil {
				return nil, fuego.HTTPError{
					Title:  "error getting request body",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			err = publisherManager.Publish(spanCtx, eventprocessing.PublishRequest{
				Topic:  "micartapro.events",
				Source: "micartapro.api.menu.interaction",
				Event:  body,
			})
			if err != nil {
				return nil, fuego.HTTPError{
					Title:  "error publishing event",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			obs.Logger.InfoContext(spanCtx, "menuInteractionRequest", "requestBody", body)
			return http.StatusOK, nil
		}, option.Summary("menuInteractionRequest"))
}
