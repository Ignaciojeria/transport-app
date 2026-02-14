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

func init() {
	ioc.Registry(
		closeJourneyHandler,
		httpserver.New,
		observability.NewObservability,
		supabaserepo.NewGetActiveJourneyByMenuID,
		supabaserepo.NewCloseJourney,
		supabaserepo.NewGetMenuIdByUserId,
		apimiddleware.NewJWTAuthMiddleware,
	)
}

func closeJourneyHandler(
	s httpserver.Server,
	obs observability.Observability,
	getActiveJourney supabaserepo.GetActiveJourneyByMenuID,
	closeJourney supabaserepo.CloseJourney,
	getMenuIdByUserId supabaserepo.GetMenuIdByUserId,
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

			if err := closeJourney(spanCtx, menuID, active.ID); err != nil {
				obs.Logger.ErrorContext(spanCtx, "error closing journey", "error", err, "menuID", menuID, "journeyID", active.ID)
				return struct{}{}, fuego.HTTPError{
					Title:  "error closing journey",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
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
