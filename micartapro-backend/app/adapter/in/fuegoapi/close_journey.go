package fuegoapi

import (
	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/adapter/out/storage"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"micartapro/app/usecase/journey"
	"micartapro/app/usecase/order"
	"net/http"
	"time"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

// CloseJourneyRequest body para cerrar jornada. SIEMPRE se requiere elegir la acción con órdenes pendientes.
type CloseJourneyRequest struct {
	// PendingOrdersAction: "cancel" = cancelar órdenes pendientes; "keep" = mantener para próxima jornada
	PendingOrdersAction string `json:"pendingOrdersAction"`
}

func init() {
	ioc.Register(closeJourneyHandler)
}

func closeJourneyHandler(
	s httpserver.Server,
	obs observability.Observability,
	getActiveJourney supabaserepo.GetActiveJourneyByMenuID,
	closeJourney supabaserepo.CloseJourney,
	userHasMenu supabaserepo.UserHasMenu,
	getOrderItems supabaserepo.GetOrderItemsForJourney,
	getJourneyStats supabaserepo.GetJourneyStats,
	releaseOrders supabaserepo.ReleaseOrdersFromJourney,
	uploadReport storage.UploadJourneyReport,
	updateReportURL supabaserepo.UpdateJourneyReportURL,
	cancelOrder order.CancelOrder,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware,
) {
	fuego.Post(s.Manager, "/api/menus/{menuId}/journeys/close",
		func(c fuego.ContextWithBody[CloseJourneyRequest]) (struct{}, error) {
			ctx := c.Context()
			spanCtx, span := obs.Tracer.Start(ctx, "closeJourney")
			defer span.End()

			menuID := c.PathParam("menuId")
			if menuID == "" {
				return struct{}{}, fuego.HTTPError{
					Title:  "menuId is required",
					Detail: "menuId path parameter is required",
					Status: http.StatusBadRequest,
				}
			}

			body, err := c.Body()
			if err != nil {
				return struct{}{}, fuego.HTTPError{
					Title:  "invalid body",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}
			if body.PendingOrdersAction != "cancel" && body.PendingOrdersAction != "keep" {
				return struct{}{}, fuego.HTTPError{
					Title:  "pendingOrdersAction required",
					Detail: "pendingOrdersAction must be 'cancel' or 'keep'",
					Status: http.StatusBadRequest,
				}
			}

			userID, ok := sharedcontext.UserIDFromContext(spanCtx)
			if !ok || userID == "" {
				return struct{}{}, fuego.HTTPError{
					Title:  "unauthorized",
					Detail: "user id not found in context",
					Status: http.StatusUnauthorized,
				}
			}

			hasMenu, err := userHasMenu(spanCtx, userID, menuID)
			if err != nil || !hasMenu {
				return struct{}{}, fuego.HTTPError{
					Title:  "menu not found",
					Detail: "menu not found or you do not own it",
					Status: http.StatusNotFound,
				}
			}

			active, err := getActiveJourney(spanCtx, menuID)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error getting active journey", "error", err, "menuID", menuID)
				return struct{}{}, fuego.HTTPError{
					Title:  "error getting journey",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			if active == nil {
				return struct{}{}, fuego.HTTPError{
					Title:  "no open journey",
					Detail: "no hay jornada abierta para este menú",
					Status: http.StatusNotFound,
				}
			}

			// Procesar órdenes pendientes según la acción elegida (SIEMPRE preguntar, nunca automático)
			items, err := getOrderItems(spanCtx, active.ID)
			if err != nil {
				obs.Logger.WarnContext(spanCtx, "error getting order items for pending action", "error", err, "journeyID", active.ID)
			}
			pendingAggregateIDs := make(map[int64]bool)
			for _, it := range items {
				if it.Status == "PENDING" || it.Status == "IN_PROGRESS" || it.Status == "READY" {
					pendingAggregateIDs[it.AggregateID] = true
				}
			}
			if err == nil {
				if body.PendingOrdersAction == "cancel" {
					for aggID := range pendingAggregateIDs {
						req := events.OrderCancelledRequest{
							AggregateID: aggID,
							Reason:      "Jornada cerrada",
							UpdatedAt:   time.Now().UTC().Format(time.RFC3339Nano),
						}
						if err := cancelOrder(spanCtx, aggID, req); err != nil {
							obs.Logger.ErrorContext(spanCtx, "error cancelling order on journey close", "error", err, "aggregateID", aggID)
							return struct{}{}, fuego.HTTPError{
								Title:  "error cancelling pending orders",
								Detail: err.Error(),
								Status: http.StatusInternalServerError,
							}
						}
					}
				} else {
					if err := releaseOrders(spanCtx, active.ID); err != nil {
						obs.Logger.ErrorContext(spanCtx, "error releasing orders from journey", "error", err, "journeyID", active.ID)
						return struct{}{}, fuego.HTTPError{
							Title:  "error releasing orders for next journey",
							Detail: err.Error(),
							Status: http.StatusInternalServerError,
						}
					}
				}
			}

			// Snapshot de totales para la jornada (desde journey_product_stats)
			var totalsSnapshot interface{}
			if stats, err := getJourneyStats(spanCtx, active.ID); err == nil {
				topProducts := make([]map[string]interface{}, 0, len(stats.Products))
				for _, p := range stats.Products {
					topProducts = append(topProducts, map[string]interface{}{
						"name":     p.ProductName,
						"quantity": p.QuantitySold,
						"revenue":  p.TotalRevenue,
					})
				}
				totalsSnapshot = map[string]interface{}{
					"totalRevenue": stats.TotalRevenue,
					"totalOrders":  stats.TotalOrders,
					"topProducts":  topProducts,
				}
			}

			if err := closeJourney(spanCtx, menuID, active.ID, totalsSnapshot); err != nil {
				obs.Logger.ErrorContext(spanCtx, "error closing journey", "error", err, "menuID", menuID, "journeyID", active.ID)
				return struct{}{}, fuego.HTTPError{
					Title:  "error closing journey",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			// Generar reporte XLSX al cerrar (snapshot de la jornada, usa items ya obtenidos)
			closedAt := time.Now().UTC()
			if len(items) > 0 {
				reportItems := make([]journey.ReportOrderItem, 0, len(items))
				for _, it := range items {
					status := it.Status
					if body.PendingOrdersAction == "cancel" && pendingAggregateIDs[it.AggregateID] {
						status = "CANCELLED"
					}
					reportItems = append(reportItems, journey.ReportOrderItem{
						OrderNumber:   it.OrderNumber,
						ItemName:      it.ItemName,
						Quantity:      it.Quantity,
						Unit:          it.Unit,
						Station:       it.Station,
						Fulfillment:   it.Fulfillment,
						Status:        status,
						RequestedTime: it.RequestedTime,
						CreatedAt:     it.CreatedAt,
						TotalPrice:    it.TotalPrice,
						TotalCost:     it.TotalCost,
					})
				}
				xlsxBytes, err := journey.GenerateJourneyReportXLSX(reportItems, active.OpenedAt, closedAt)
				if err != nil {
					obs.Logger.WarnContext(spanCtx, "error generating report xlsx", "error", err, "journeyID", active.ID)
				} else {
					reportURL, err := uploadReport(spanCtx, active.ID, xlsxBytes)
					if err != nil {
						obs.Logger.WarnContext(spanCtx, "error uploading report", "error", err, "journeyID", active.ID)
					} else if err := updateReportURL(spanCtx, active.ID, reportURL); err != nil {
						obs.Logger.WarnContext(spanCtx, "error updating journey report url", "error", err, "journeyID", active.ID)
					} else {
						obs.Logger.InfoContext(spanCtx, "journey report generated", "journeyID", active.ID)
					}
				}
			}

			obs.Logger.InfoContext(spanCtx, "journey closed", "menuID", menuID, "journeyID", active.ID)
			return struct{}{}, nil
		},
		option.Summary("Close the active journey for the menu"),
		option.Description("Sets the active (OPEN) journey to CLOSED and sets closed_at. Returns 404 if there is no open journey for this menu."),
		option.Tags("menu", "journey"),
		option.Middleware(jwtAuthMiddleware),
	)
}
