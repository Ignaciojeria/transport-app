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
	"github.com/google/uuid"
)

func init() {
	ioc.Registry(
		createTenant,
		httpserver.New,
		gcppublisher.NewApplicationEvents,
		observability.NewObservability,
		usecase.NewCreateTenant)
}
func createTenant(
	s httpserver.Server,
	publish gcppublisher.ApplicationEvents,
	obs observability.Observability,
	createOrg usecase.CreateTenant) {
	fuego.Post(s.Manager, "/tenants",
		func(c fuego.ContextWithBody[request.CreateTenantRequest]) (response.CreateTenantResponse, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "createTenant")
			defer span.End()
			requestBody, err := c.Body()
			if err != nil {
				return response.CreateTenantResponse{}, err
			}
			requestBody.ID = uuid.NewString()
			eventPayload, _ := json.Marshal(requestBody)
			eventCtx := sharedcontext.AddEventContextToBaggage(spanCtx,
				sharedcontext.EventContext{
					EntityType: "tenant",
					EventType:  "tenantSubmitted",
				})

			if err := publish(eventCtx, domain.Outbox{
				Payload: eventPayload,
			}); err != nil {
				return response.CreateTenantResponse{}, fuego.HTTPError{
					Title:  "error creating tenant",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			obs.Logger.InfoContext(spanCtx,
				"TENANT_SUBMISSION_SUCCEEDED",
				slog.Any("payload", requestBody))

			return response.CreateTenantResponse{
				ID:      requestBody.ID,
				Country: requestBody.Country,
				Tenant:  requestBody.ID + "-" + requestBody.Country,
				Message: "Tenant submitted successfully",
			}, nil
		},
		option.Summary("create tenant"),
		option.Tags(tagTenants),
	)
}
