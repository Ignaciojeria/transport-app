package fuegoapi

import (
	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"net/http"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

// JourneyStatsResponse es la respuesta del endpoint de estadísticas.
type JourneyStatsResponse struct {
	TotalRevenue    float64               `json:"totalRevenue"`
	TotalCost       float64               `json:"totalCost"`
	TotalOrders     int                   `json:"totalOrders"`
	ItemsOrdered    int                   `json:"itemsOrdered"`
	Products        []ProductStatResponse `json:"products"`
	RevenueByStatus RevenueByStatusResp   `json:"revenueByStatus"`
	OrdersByStatus  OrdersByStatusResp    `json:"ordersByStatus"`
}

// RevenueByStatusResp ventas por estado.
type RevenueByStatusResp struct {
	Delivered  float64 `json:"delivered"`
	Dispatched float64 `json:"dispatched"`
	Pending    float64 `json:"pending"`
	Cancelled  float64 `json:"cancelled"`
}

// OrdersByStatusResp órdenes por estado.
type OrdersByStatusResp struct {
	Delivered  int `json:"delivered"`
	Dispatched int `json:"dispatched"`
	Pending    int `json:"pending"`
	Cancelled  int `json:"cancelled"`
}

// ProductStatResponse es un producto con sus estadísticas.
type ProductStatResponse struct {
	ProductName          string  `json:"productName"`
	QuantitySold         int     `json:"quantitySold"`
	TotalRevenue         float64 `json:"totalRevenue"`
	TotalCost            float64 `json:"totalCost"`
	Percentage           float64 `json:"percentage"`
	PercentageByQuantity float64 `json:"percentageByQuantity"`
}

func init() {
	ioc.Register(getJourneyStatsHandler)
}

func getJourneyStatsHandler(
	s httpserver.Server,
	obs observability.Observability,
	getJourneyStats supabaserepo.GetJourneyStats,
	userHasMenu supabaserepo.UserHasMenu,
	getJourneys supabaserepo.GetJourneysByMenuID,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware,
) {
	fuego.Get(s.Manager, "/api/menus/{menuId}/journeys/{journeyId}/stats",
		func(c fuego.ContextNoBody) (JourneyStatsResponse, error) {
			ctx := c.Context()
			spanCtx, span := obs.Tracer.Start(ctx, "getJourneyStats")
			defer span.End()

			menuID := c.PathParam("menuId")
			journeyID := c.PathParam("journeyId")
			if menuID == "" || journeyID == "" {
				return JourneyStatsResponse{}, fuego.HTTPError{
					Title:  "menuId and journeyId are required",
					Detail: "path parameters are required",
					Status: http.StatusBadRequest,
				}
			}

			userID, ok := sharedcontext.UserIDFromContext(spanCtx)
			if !ok || userID == "" {
				return JourneyStatsResponse{}, fuego.HTTPError{
					Title:  "unauthorized",
					Detail: "user id not found in context",
					Status: http.StatusUnauthorized,
				}
			}

			hasMenu, err := userHasMenu(spanCtx, userID, menuID)
			if err != nil || !hasMenu {
				return JourneyStatsResponse{}, fuego.HTTPError{
					Title:  "menu not found",
					Detail: "menu not found or you do not own it",
					Status: http.StatusNotFound,
				}
			}

			// Verificar que la jornada pertenece al menú
			journeys, err := getJourneys(spanCtx, menuID, 100)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error listing journeys", "error", err, "menuID", menuID)
				return JourneyStatsResponse{}, fuego.HTTPError{
					Title:  "error getting journey",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			var found bool
			for _, j := range journeys {
				if j.ID == journeyID {
					found = true
					break
				}
			}
			if !found {
				return JourneyStatsResponse{}, fuego.HTTPError{
					Title:  "journey not found",
					Detail: "journey not found or does not belong to this menu",
					Status: http.StatusNotFound,
				}
			}

			stats, err := getJourneyStats(spanCtx, journeyID)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error getting journey stats", "error", err, "journeyID", journeyID)
				return JourneyStatsResponse{}, fuego.HTTPError{
					Title:  "error getting stats",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			products := make([]ProductStatResponse, 0, len(stats.Products))
			for _, p := range stats.Products {
				products = append(products, ProductStatResponse{
					ProductName:          p.ProductName,
					QuantitySold:         p.QuantitySold,
					TotalRevenue:         p.TotalRevenue,
					TotalCost:            p.TotalCost,
					Percentage:           p.Percentage,
					PercentageByQuantity: p.PercentageByQuantity,
				})
			}

			return JourneyStatsResponse{
				TotalRevenue: stats.TotalRevenue,
				TotalCost:    stats.TotalCost,
				TotalOrders:  stats.TotalOrders,
				ItemsOrdered: stats.ItemsOrdered,
				Products:     products,
				RevenueByStatus: RevenueByStatusResp{
					Delivered:  stats.RevenueByStatus.Delivered,
					Dispatched: stats.RevenueByStatus.Dispatched,
					Pending:    stats.RevenueByStatus.Pending,
					Cancelled:  stats.RevenueByStatus.Cancelled,
				},
				OrdersByStatus: OrdersByStatusResp{
					Delivered:  stats.OrdersByStatus.Delivered,
					Dispatched: stats.OrdersByStatus.Dispatched,
					Pending:    stats.OrdersByStatus.Pending,
					Cancelled:  stats.OrdersByStatus.Cancelled,
				},
			}, nil
		},
		option.Summary("Get journey statistics"),
		option.Description("Returns sales statistics for a journey: total revenue, order count, and product breakdown with percentages. Suitable for pie charts and analytics."),
		option.Tags("menu", "journey", "stats"),
		option.Middleware(jwtAuthMiddleware),
	)
}
