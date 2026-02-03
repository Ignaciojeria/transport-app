package fuegoapi

import (
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/usecase/order"
	"net/http"
	"strconv"
	"time"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func init() {
	ioc.Registry(
		startOrderPreparation,
		httpserver.New,
		observability.NewObservability,
		order.NewStartPreparation,
	)
}

func startOrderPreparation(
	s httpserver.Server,
	obs observability.Observability,
	startPreparationUseCase order.StartPreparation,
) {
	fuego.Post(s.Manager, "/menu/{menuId}/orders/{aggregateId}/start-preparation",
		func(c fuego.ContextWithBody[events.OrderStartedPreparationRequest]) (map[string]string, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "startOrderPreparation")
			defer span.End()

			// Obtener el aggregateId del parámetro de ruta
			aggregateIDStr := c.PathParam("aggregateId")
			aggregateID, err := strconv.ParseInt(aggregateIDStr, 10, 64)
			if err != nil {
				return nil, fuego.HTTPError{
					Title:  "invalid aggregateId",
					Detail: "aggregateId must be a valid integer",
					Status: http.StatusBadRequest,
				}
			}

			// Obtener el body del request
			request, err := c.Body()
			if err != nil {
				return nil, fuego.HTTPError{
					Title:  "error getting request body",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}

			// Validar que station esté presente
			if request.Station == "" {
				return nil, fuego.HTTPError{
					Title:  "station is required",
					Detail: "station field is required (KITCHEN or BAR)",
					Status: http.StatusBadRequest,
				}
			}

			// Fijar updatedAt y aggregateID
			request.UpdatedAt = time.Now().Format(time.RFC3339Nano)
			request.AggregateID = aggregateID

			// Llamar al caso de uso
			err = startPreparationUseCase(spanCtx, aggregateID, request)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error starting preparation", "error", err, "aggregateID", aggregateID)
				return nil, fuego.HTTPError{
					Title:  "error starting preparation",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			obs.Logger.InfoContext(spanCtx, "preparation started successfully", "aggregateID", aggregateID)
			return map[string]string{"status": "success"}, nil
		},
		option.Summary("startOrderPreparation"),
		option.Tags("orders"),
		option.Path("menuId", "string", param.Required()),
		option.Path("aggregateId", "string", param.Required()))
}
