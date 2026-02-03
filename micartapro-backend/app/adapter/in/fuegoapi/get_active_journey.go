package fuegoapi

import (
	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"micartapro/app/usecase/journey"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(
		getActiveJourneyHandler,
		httpserver.New,
		observability.NewObservability,
		supabaserepo.NewGetActiveJourneyByMenuID,
		supabaserepo.NewGetMenuIdByUserId,
		apimiddleware.NewJWTAuthMiddleware,
	)
}

func getActiveJourneyHandler(
	s httpserver.Server,
	obs observability.Observability,
	getActiveJourney supabaserepo.GetActiveJourneyByMenuID,
	getMenuIdByUserId supabaserepo.GetMenuIdByUserId,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware,
) {
	fuego.Get(s.Manager, "/api/menus/{menuId}/active-journey",
		func(c fuego.ContextNoBody) (*journey.Journey, error) {
			ctx := c.Context()
			spanCtx, span := obs.Tracer.Start(ctx, "getActiveJourney")
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

			j, err := getActiveJourney(spanCtx, menuID)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error getting active journey", "error", err, "menuID", menuID)
				return nil, fuego.HTTPError{
					Title:  "error getting journey",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			if j == nil {
				return nil, fuego.HTTPError{
					Title:  "no active journey",
					Detail: "there is no open journey for this menu",
					Status: http.StatusNotFound,
				}
			}
			return j, nil
		},
		option.Summary("Get active journey for menu"),
		option.Description("Returns the currently open journey for the given menu. Requires the authenticated user to own the menu."),
		option.Tags("menu", "journey"),
		option.Middleware(jwtAuthMiddleware),
	)
}
