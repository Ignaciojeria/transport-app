package fuegoapi

import (
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/usecase/order"
	"net/http"
	"strconv"
	"time"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func init() {
	ioc.Register(dispatchOrder)
}

func dispatchOrder(
	s httpserver.Server,
	obs observability.Observability,
	dispatchOrderUseCase order.DispatchOrder,
) {
	fuego.Post(s.Manager, "/menu/{menuId}/orders/{aggregateId}/dispatch",
		func(c fuego.ContextWithBody[events.OrderDeliveredRequest]) (map[string]string, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "dispatchOrder")
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

			// Obtener el body del request (puede estar vacío)
			request, err := c.Body()
			if err != nil {
				return nil, fuego.HTTPError{
					Title:  "error getting request body",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}

			// Conservar updatedAt UTC del frontend; si no viene, fijar en servidor (UTC)
			if request.UpdatedAt == "" {
				request.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
			}
			request.AggregateID = aggregateID

			// Llamar al caso de uso
			err = dispatchOrderUseCase(spanCtx, aggregateID, request)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error dispatching order", "error", err, "aggregateID", aggregateID)
				return nil, fuego.HTTPError{
					Title:  "error dispatching order",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			obs.Logger.InfoContext(spanCtx, "order delivered successfully", "aggregateID", aggregateID)
			return map[string]string{"status": "success"}, nil
		},
		option.Summary("dispatchOrder"),
		option.Tags("orders"),
		option.Path("menuId", "string", param.Required()),
		option.Path("aggregateId", "string", param.Required()))
}
