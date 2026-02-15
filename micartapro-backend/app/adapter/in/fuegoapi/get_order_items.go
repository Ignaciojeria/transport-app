package fuegoapi

import (
	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"net/http"
	"strconv"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

// OrderItemResponse representa un ítem de orden en la respuesta API.
type OrderItemResponse struct {
	ItemName   string   `json:"itemName"`
	Quantity   int      `json:"quantity"`
	Unit       string   `json:"unit"`
	TotalPrice float64  `json:"totalPrice"`
	Station    *string  `json:"station,omitempty"`
}

func init() {
	ioc.Registry(
		getOrderItemsHandler,
		httpserver.New,
		observability.NewObservability,
		supabaserepo.NewGetOrderItemsByAggregateID,
		supabaserepo.NewGetMenuIdByUserId,
		apimiddleware.NewJWTAuthMiddleware,
	)
}

func getOrderItemsHandler(
	s httpserver.Server,
	obs observability.Observability,
	getOrderItems supabaserepo.GetOrderItemsByAggregateID,
	getMenuIdByUserId supabaserepo.GetMenuIdByUserId,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware,
) {
	fuego.Get(s.Manager, "/api/menus/{menuId}/orders/{aggregateId}/items",
		func(c fuego.ContextNoBody) ([]OrderItemResponse, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "getOrderItems")
			defer span.End()

			menuID := c.PathParam("menuId")
			if menuID == "" {
				return nil, fuego.HTTPError{
					Title:  "menuId is required",
					Detail: "menuId path parameter is required",
					Status: http.StatusBadRequest,
				}
			}

			aggregateIDStr := c.PathParam("aggregateId")
			aggregateID, err := strconv.ParseInt(aggregateIDStr, 10, 64)
			if err != nil {
				return nil, fuego.HTTPError{
					Title:  "invalid aggregateId",
					Detail: "aggregateId must be a valid integer",
					Status: http.StatusBadRequest,
				}
			}

			userID, ok := sharedcontext.UserIDFromContext(spanCtx)
			if !ok || userID == "" {
				return nil, fuego.HTTPError{
					Title:  "unauthorized",
					Detail: "user id not found in context",
					Status: http.StatusUnauthorized,
				}
			}

			userMenuID, err := getMenuIdByUserId(spanCtx, userID)
			if err != nil || userMenuID != menuID {
				return nil, fuego.HTTPError{
					Title:  "menu not found",
					Detail: "menu not found or you do not own it",
					Status: http.StatusNotFound,
				}
			}

			items, err := getOrderItems(spanCtx, menuID, aggregateID)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error getting order items", "error", err, "aggregateID", aggregateID)
				return nil, fuego.HTTPError{
					Title:  "error getting order items",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			out := make([]OrderItemResponse, 0, len(items))
			for _, it := range items {
				out = append(out, OrderItemResponse{
					ItemName:   it.ItemName,
					Quantity:   it.Quantity,
					Unit:       it.Unit,
					TotalPrice: it.TotalPrice,
					Station:    it.Station,
				})
			}
			return out, nil
		},
		option.Summary("Get order items"),
		option.Description("Returns items for an order by aggregate_id."),
		option.Tags("menu", "orders"),
		option.Path("menuId", "string", param.Required()),
		option.Path("aggregateId", "string", param.Required()),
		option.Middleware(jwtAuthMiddleware),
	)
}
