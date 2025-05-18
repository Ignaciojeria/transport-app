package fuegoapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/adapter/out/gcppublisher"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"
	"transport-app/app/usecase"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/baggage"
)

func init() {
	ioc.Registry(
		createTenant,
		httpserver.New,
		gcppublisher.NewApplicationEvents,
		observability.NewObservability,
		usecase.NewCreateTenant,
		tidbrepository.NewFindAccountByEmail,
	)
}

func createTenant(
	s httpserver.Server,
	publish gcppublisher.ApplicationEvents,
	obs observability.Observability,
	createOrg usecase.CreateTenant,
	findAccount tidbrepository.FindAccountByEmail,
) {
	fuego.Post(s.Manager, "/tenants",
		func(c fuego.ContextWithBody[request.CreateTenantRequest]) (response.CreateTenantResponse, error) {
			requestBody, err := c.Body()
			if err != nil {
				return response.CreateTenantResponse{}, err
			}

			// Verificar si el email existe
			account, err := findAccount(c.Context(), requestBody.Email)
			if err != nil {
				return response.CreateTenantResponse{}, fuego.HTTPError{
					Title:  "error finding account",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			if account.Email == "" {
				return response.CreateTenantResponse{}, fuego.HTTPError{
					Title:  "email not registered",
					Detail: "the email is not registered in the system",
					Status: http.StatusBadRequest,
				}
			}

			// Asignar nuevo UUID como tenant ID
			requestBody.ID = uuid.NewString()

			// Crear baggage manualmente con tenant y country antes de iniciar la traza
			mTenant, _ := baggage.NewMember(sharedcontext.BaggageTenantID, requestBody.ID)
			mCountry, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, requestBody.Country)
			bag, _ := baggage.New(mTenant, mCountry)
			baggageCtx := baggage.ContextWithBaggage(c.Context(), bag)

			// Iniciar traza con el contexto que ya tiene el baggage
			spanCtx, span := obs.Tracer.Start(baggageCtx, "createTenant")
			defer span.End()

			// Enriquecer el contexto con metadata de evento
			eventCtx := sharedcontext.AddEventContextToBaggage(spanCtx,
				sharedcontext.EventContext{
					EntityType: "tenant",
					EventType:  "tenantSubmitted",
				})

			eventPayload, _ := json.Marshal(requestBody)

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
