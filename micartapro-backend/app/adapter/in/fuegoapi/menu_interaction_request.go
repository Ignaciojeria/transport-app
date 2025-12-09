package fuegoapi

import (
	"micartapro/app/adapter/out/publisher"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"net/http"
	"time"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

type ChatMessage struct {
	Role    string `json:"role"` // "user" | "assistant" | "system"
	Content string `json:"content"`
}

type MenuPreferences struct {
	Language string `json:"language"`
}

type MenuInteractionRequest struct {
	Message         string          `json:"message"`
	History         []ChatMessage   `json:"history,omitempty"`
	MenuPreferences MenuPreferences `json:"menuPreferences,omitempty"`
}

func (m MenuInteractionRequest) ToCloudEvent() cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("menu.interaction.request")
	event.SetType("menu.interaction.requested")
	event.SetSource("micartapro.api.menu.interaction")
	event.SetData(cloudevents.ApplicationJSON, m)
	event.SetTime(time.Now())
	return event
}

func init() {
	ioc.Registry(
		menuInteractionRequest,
		httpserver.New,
		observability.NewObservability,
		publisher.NewPublishEvents)
}
func menuInteractionRequest(
	s httpserver.Server,
	obs observability.Observability,
	publish publisher.PublishEvents) {
	fuego.Post(s.Manager, "/menu/interaction",
		func(c fuego.ContextWithBody[MenuInteractionRequest]) (any, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "menuGenerationRequest")
			defer span.End()
			body, err := c.Body()
			if err != nil {
				return nil, fuego.HTTPError{
					Title:  "error al obtener el body",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			err = publish(spanCtx, body)
			if err != nil {
				return nil, fuego.HTTPError{
					Title:  "error publishing event",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			obs.Logger.InfoContext(spanCtx, "menuGenerationRequest", "requestBody", body)
			return http.StatusOK, nil
		}, option.Summary("menuGenerationRequest"))
}
