package fuegoapi

import (
	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/adapter/out/storage"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"micartapro/app/usecase/journey"
	"net/http"
	"time"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(
		closeJourneyHandler,
		httpserver.New,
		observability.NewObservability,
		supabaserepo.NewGetActiveJourneyByMenuID,
		supabaserepo.NewCloseJourney,
		supabaserepo.NewGetMenuIdByUserId,
		supabaserepo.NewGetOrderItemsForJourney,
		supabaserepo.NewGetJourneyStats,
		storage.NewUploadJourneyReport,
		supabaserepo.NewUpdateJourneyReportURL,
		apimiddleware.NewJWTAuthMiddleware,
	)
}

func closeJourneyHandler(
	s httpserver.Server,
	obs observability.Observability,
	getActiveJourney supabaserepo.GetActiveJourneyByMenuID,
	closeJourney supabaserepo.CloseJourney,
	getMenuIdByUserId supabaserepo.GetMenuIdByUserId,
	getOrderItems supabaserepo.GetOrderItemsForJourney,
	getJourneyStats supabaserepo.GetJourneyStats,
	uploadReport storage.UploadJourneyReport,
	updateReportURL supabaserepo.UpdateJourneyReportURL,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware,
) {
	fuego.Post(s.Manager, "/api/menus/{menuId}/journeys/close",
		func(c fuego.ContextNoBody) (struct{}, error) {
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

			userID, ok := sharedcontext.UserIDFromContext(spanCtx)
			if !ok || userID == "" {
				return struct{}{}, fuego.HTTPError{
					Title:  "unauthorized",
					Detail: "user id not found in context",
					Status: http.StatusUnauthorized,
				}
			}

			userMenuID, err := getMenuIdByUserId(spanCtx, userID)
			if err != nil || userMenuID != menuID {
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

			// Generar reporte XLSX al cerrar (snapshot de la jornada)
			closedAt := time.Now().UTC()
			items, err := getOrderItems(spanCtx, active.ID)
			if err != nil {
				obs.Logger.WarnContext(spanCtx, "error getting order items for report", "error", err, "journeyID", active.ID)
			} else {
				reportItems := make([]journey.ReportOrderItem, 0, len(items))
				for _, it := range items {
					reportItems = append(reportItems, journey.ReportOrderItem{
						OrderNumber:   it.OrderNumber,
						ItemName:      it.ItemName,
						Quantity:      it.Quantity,
						Unit:          it.Unit,
						Station:       it.Station,
						Fulfillment:   it.Fulfillment,
						Status:        it.Status,
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
