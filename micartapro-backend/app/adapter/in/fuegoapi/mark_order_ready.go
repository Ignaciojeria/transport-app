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
		markOrderReady,
		httpserver.New,
		observability.NewObservability,
		order.NewMarkReady,
	)
}

func markOrderReady(
	s httpserver.Server,
	obs observability.Observability,
	markReadyUseCase order.MarkReady,
) {
	fuego.Post(s.Manager, "/menu/{menuId}/orders/{aggregateId}/mark-ready",
		func(c fuego.ContextWithBody[events.OrderItemReadyRequest]) (map[string]string, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "markOrderReady")
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

			// Validar que itemKeys esté presente y no vacío
			if len(request.ItemKeys) == 0 {
				return nil, fuego.HTTPError{
					Title:  "itemKeys is required",
					Detail: "itemKeys field is required and must contain at least one item key",
					Status: http.StatusBadRequest,
				}
			}

			// Fijar updatedAt y aggregateID
			request.UpdatedAt = time.Now().Format(time.RFC3339Nano)
			request.AggregateID = aggregateID

			// Llamar al caso de uso
			err = markReadyUseCase(spanCtx, aggregateID, request)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error marking ready", "error", err, "aggregateID", aggregateID)
				return nil, fuego.HTTPError{
					Title:  "error marking ready",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			obs.Logger.InfoContext(spanCtx, "items marked ready successfully", "aggregateID", aggregateID)
			return map[string]string{"status": "success"}, nil
		},
		option.Summary("markOrderReady"),
		option.Tags("orders"),
		option.Path("menuId", "string", param.Required()),
		option.Path("aggregateId", "string", param.Required()))
}
