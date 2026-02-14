package fuegoapi

import (
	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

// JourneyListItemResponse es la respuesta de una jornada en el listado.
type JourneyListItemResponse struct {
	ID            string   `json:"id"`
	MenuID        string   `json:"menuId"`
	Status        string   `json:"status"`
	OpenedAt      string   `json:"openedAt"`
	ClosedAt      *string  `json:"closedAt,omitempty"`
	ReportPDFURL  *string  `json:"reportPdfUrl,omitempty"`
	ReportXLSXURL *string  `json:"reportXlsxUrl,omitempty"`
}

func init() {
	ioc.Registry(
		listJourneysHandler,
		httpserver.New,
		observability.NewObservability,
		supabaserepo.NewGetJourneysByMenuID,
		supabaserepo.NewGetMenuIdByUserId,
		apimiddleware.NewJWTAuthMiddleware,
	)
}

func listJourneysHandler(
	s httpserver.Server,
	obs observability.Observability,
	getJourneys supabaserepo.GetJourneysByMenuID,
	getMenuIdByUserId supabaserepo.GetMenuIdByUserId,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware,
) {
	fuego.Get(s.Manager, "/api/menus/{menuId}/journeys",
		func(c fuego.ContextNoBody) ([]JourneyListItemResponse, error) {
			ctx := c.Context()
			spanCtx, span := obs.Tracer.Start(ctx, "listJourneys")
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

			items, err := getJourneys(spanCtx, menuID, 50)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error listing journeys", "error", err, "menuID", menuID)
				return nil, fuego.HTTPError{
					Title:  "error listing journeys",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			out := make([]JourneyListItemResponse, 0, len(items))
			for _, j := range items {
				closedAt := ""
				if j.ClosedAt != nil {
					closedAt = j.ClosedAt.Format("2006-01-02T15:04:05Z07:00")
				}
				resp := JourneyListItemResponse{
					ID:       j.ID,
					MenuID:   j.MenuID,
					Status:   j.Status,
					OpenedAt: j.OpenedAt.Format("2006-01-02T15:04:05Z07:00"),
					ReportPDFURL:  j.ReportPDFURL,
					ReportXLSXURL: j.ReportXLSXURL,
				}
				if j.ClosedAt != nil {
					resp.ClosedAt = &closedAt
				}
				out = append(out, resp)
			}
			return out, nil
		},
		option.Summary("List journeys for the menu"),
		option.Description("Returns journeys (open and closed) with report URLs when available. Ordered by opened_at desc."),
		option.Tags("menu", "journey"),
		option.Middleware(jwtAuthMiddleware),
	)
}
