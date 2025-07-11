package fuegoapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/adapter/out/natspublisher"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	/*
		ioc.Registry(generateTestData,
			httpserver.New,
			natspublisher.NewApplicationEvents,
			observability.NewObservability)
	*/
	ioc.Registry(
		optimizeFleet,
		httpserver.New,
		natspublisher.NewApplicationEvents,
		observability.NewObservability)
}

func optimizeFleet(
	s httpserver.Server,
	publish natspublisher.ApplicationEvents,
	obs observability.Observability) {
	fuego.Post(s.Manager, "/optimize/fleet",
		func(c fuego.ContextWithBody[request.OptimizeFleetRequest]) (response.OptimizationResponse, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "optimization")
			defer span.End()

			spanCtx = sharedcontext.WithAccessToken(spanCtx, c.Header("X-Access-Token"))

			requestBody, err := c.Body()
			if err != nil {
				return response.OptimizationResponse{}, err
			}

			eventPayload, _ := json.Marshal(requestBody)

			eventCtx := sharedcontext.AddEventContextToBaggage(spanCtx,
				sharedcontext.EventContext{
					EntityType: "optimization",
					EventType:  "optimizationRequested",
				})

			if err := publish(eventCtx, domain.Outbox{
				Payload: eventPayload,
			}); err != nil {
				return response.OptimizationResponse{}, fuego.HTTPError{
					Title:  "error requesting optimization",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			obs.Logger.InfoContext(spanCtx,
				"OPTIMIZATION_REQUEST_SUBMITTED",
				slog.Any("payload", requestBody))

			return response.OptimizationResponse{
				TraceID: span.SpanContext().TraceID().String(),
			}, nil
		}, option.Summary("optimize fleet"), option.Tags("optimization"),
		option.Header("X-Access-Token", "api access token"),
	)
}

/*
// generateTestData crea un endpoint que genera datos de prueba masivos usando el test data generator
func generateTestData(
	s httpserver.Server,
	publish natspublisher.ApplicationEvents,
	obs observability.Observability) {
	fuego.Post(s.Manager, "/optimize/fleet/test-data",
		func(c fuego.ContextNoBody) (response.OptimizationResponse, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "generate-test-data")
			defer span.End()

			// Generar datos de prueba masivos usando el test data generator
			testData := request.GenerateMassiveTestData()

			eventPayload, _ := json.Marshal(testData)

			eventCtx := sharedcontext.AddEventContextToBaggage(spanCtx,
				sharedcontext.EventContext{
					EntityType: "optimization",
					EventType:  "optimizationRequested",
				})

			if err := publish(eventCtx, domain.Outbox{
				Payload: eventPayload,
			}); err != nil {
				return response.OptimizationResponse{}, fuego.HTTPError{
					Title:  "error requesting optimization with test data",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			obs.Logger.InfoContext(spanCtx,
				"TEST_DATA_OPTIMIZATION_REQUEST_SUBMITTED",
				slog.String("planReferenceID", testData.PlanReferenceID),
				slog.Int("vehiclesCount", len(testData.Vehicles)),
				slog.Int("visitsCount", len(testData.Visits)))

			return response.OptimizationResponse{
				TraceID: span.SpanContext().TraceID().String(),
			}, nil
		}, option.Summary("generate test data for fleet optimization"), option.Tags("testing"))
}
*/
