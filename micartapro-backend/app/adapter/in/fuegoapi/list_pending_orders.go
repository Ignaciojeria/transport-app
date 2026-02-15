package fuegoapi

import (
	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"net/http"
	"time"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

// PendingOrderItemResponse representa un ítem de orden en la respuesta.
type PendingOrderItemResponse struct {
	ItemName   string  `json:"itemName"`
	Quantity   int     `json:"quantity"`
	Unit       string  `json:"unit"`
	TotalPrice float64 `json:"totalPrice"`
}

// PendingOrderResponse representa una orden pendiente en la respuesta API.
type PendingOrderResponse struct {
	AggregateID  int64                         `json:"aggregateId"`
	TrackingID   string                        `json:"trackingId"`
	CreatedAt    string                        `json:"createdAt"`
	ScheduledFor *string                       `json:"scheduledFor,omitempty"`
	TotalAmount  int64                         `json:"totalAmount"`
	Status       string                        `json:"status"`
	Items        []PendingOrderItemResponse    `json:"items"`
}

func init() {
	ioc.Registry(
		listPendingOrdersHandler,
		httpserver.New,
		observability.NewObservability,
		supabaserepo.NewGetPendingOrders,
		supabaserepo.NewGetMenuIdByUserId,
		apimiddleware.NewJWTAuthMiddleware,
	)
}

func listPendingOrdersHandler(
	s httpserver.Server,
	obs observability.Observability,
	getPendingOrders supabaserepo.GetPendingOrders,
	getMenuIdByUserId supabaserepo.GetMenuIdByUserId,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware,
) {
	fuego.Get(s.Manager, "/api/menus/{menuId}/pending-orders",
		func(c fuego.ContextNoBody) ([]PendingOrderResponse, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "listPendingOrders")
			defer span.End()

			menuID := c.PathParam("menuId")
			if menuID == "" {
				return nil, fuego.HTTPError{
					Title:  "menuId is required",
					Detail: "menuId path parameter is required",
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

			// Filtros: fromDate, toDate (UTC ISO-8601). Parseo higiénico: validar formato sin alterar zona horaria.
			validISO8601 := func(s string) bool {
				_, err := time.Parse(time.RFC3339Nano, s)
				if err == nil {
					return true
				}
				_, err = time.Parse(time.RFC3339, s)
				return err == nil
			}
			var filter *supabaserepo.GetPendingOrdersFilter
			if fromStr := c.QueryParam("fromDate"); fromStr != "" && validISO8601(fromStr) {
				if filter == nil {
					filter = &supabaserepo.GetPendingOrdersFilter{}
				}
				filter.FromDate = &fromStr
			}
			if toStr := c.QueryParam("toDate"); toStr != "" && validISO8601(toStr) {
				if filter == nil {
					filter = &supabaserepo.GetPendingOrdersFilter{}
				}
				filter.ToDate = &toStr
			}

			orders, err := getPendingOrders(spanCtx, menuID, filter)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error listing pending orders", "error", err, "menuID", menuID)
				return nil, fuego.HTTPError{
					Title:  "error listing pending orders",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			out := make([]PendingOrderResponse, 0, len(orders))
			for _, o := range orders {
				items := make([]PendingOrderItemResponse, 0, len(o.Items))
				for _, it := range o.Items {
					items = append(items, PendingOrderItemResponse{
						ItemName:   it.ItemName,
						Quantity:   it.Quantity,
						Unit:       it.Unit,
						TotalPrice: it.TotalPrice,
					})
				}
				out = append(out, PendingOrderResponse{
					AggregateID:  o.AggregateID,
					TrackingID:   o.TrackingID,
					CreatedAt:    o.CreatedAt,
					ScheduledFor: o.ScheduledFor,
					TotalAmount:  o.TotalAmount,
					Status:       o.Status,
					Items:        items,
				})
			}
			return out, nil
		},
		option.Summary("List pending orders"),
		option.Description("Returns orders with journey_id IS NULL (not yet assigned to a journey). Optional query params: fromDate, toDate (RFC3339)."),
		option.Tags("menu", "orders"),
		option.Path("menuId", "string", param.Required()),
		option.Middleware(jwtAuthMiddleware),
	)
}
