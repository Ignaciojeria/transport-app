package fuegoapi

import (
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/usecase/order"
	"net/http"
	"time"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func init() {
	ioc.Register(createOrder)
}

func createOrder(
	s httpserver.Server,
	obs observability.Observability,
	createOrderUseCase order.CreateOrder,
) {
	fuego.Post(s.Manager, "/menu/{menuId}/orders",
		func(c fuego.ContextWithBody[events.CreateOrderRequest]) (order.CreateOrderResult, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "createOrder")
			defer span.End()

			// Obtener el menuId del parámetro de ruta
			menuID := c.PathParam("menuId")
			if menuID == "" {
				return order.CreateOrderResult{}, fuego.HTTPError{
					Title:  "menuId is required",
					Detail: "menuId parameter is required",
					Status: http.StatusBadRequest,
				}
			}

			// Obtener el body del request
			request, err := c.Body()
			if err != nil {
				return order.CreateOrderResult{}, fuego.HTTPError{
					Title:  "error getting request body",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}

			// Fijar createdAt en el momento en que llega al controlador
			request.CreatedAt = time.Now().Format(time.RFC3339Nano)

			// Llamar al caso de uso
			result, err := createOrderUseCase(spanCtx, menuID, request)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error creating order", "error", err, "menuID", menuID)
				return order.CreateOrderResult{}, fuego.HTTPError{
					Title:  "error creating order",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			obs.Logger.InfoContext(spanCtx, "order created successfully", "menuID", menuID, "orderNumber", result.OrderNumber)
			return result, nil
		},
		option.Summary("createOrder"),
		option.Tags("orders"),
		option.Path("menuId", "string", param.Required()))
}
