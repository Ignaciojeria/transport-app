package fuegoapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/adapter/out/gcppublisher"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"
	"transport-app/app/usecase"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(
		register,
		httpserver.New,
		gcppublisher.NewApplicationEvents,
		observability.NewObservability,
		usecase.NewRegister)
}
func register(
	s httpserver.Server,
	publish gcppublisher.ApplicationEvents,
	obs observability.Observability,
	register usecase.Register) {
	fuego.Post(s.Manager, "/register",
		func(c fuego.ContextWithBody[request.RegisterRequest]) (response.RegisterResponse, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "register")
			defer span.End()
			req, err := c.Body()
			if err != nil {
				return response.RegisterResponse{}, err
			}

			eventPayload, _ := json.Marshal(req)

			eventCtx := sharedcontext.AddEventContextToBaggage(spanCtx,
				sharedcontext.EventContext{
					EntityType: "registration",
					EventType:  "registrationSubmitted",
				})

			if err := publish(eventCtx, domain.Outbox{
				Payload: eventPayload,
			}); err != nil {
				return response.RegisterResponse{}, fuego.HTTPError{
					Title:  "registration error",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			obs.Logger.InfoContext(spanCtx,
				"REGISTRATION_SUBMITTED",
				slog.Any("payload", eventPayload))

			return response.RegisterResponse{
				Message: "user submitted successfully",
			}, nil
		},
		option.Tags(tagRegistration),
		option.Summary("register account"))
}
