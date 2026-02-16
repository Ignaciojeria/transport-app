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

type MenuSlugResponse struct {
	Slug string `json:"slug"`
}

func init() {
	ioc.Registry(
		getMenuSlugHandler,
		httpserver.New,
		observability.NewObservability,
		supabaserepo.NewGetActiveSlugByMenuID,
		supabaserepo.NewUserHasMenu,
		apimiddleware.NewJWTAuthMiddleware,
	)
}

func getMenuSlugHandler(
	s httpserver.Server,
	obs observability.Observability,
	getActiveSlug supabaserepo.GetActiveSlugByMenuID,
	userHasMenu supabaserepo.UserHasMenu,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware,
) {
	fuego.Get(s.Manager, "/api/menus/{menuId}/slug",
		func(c fuego.ContextNoBody) (MenuSlugResponse, error) {
			ctx := c.Context()
			spanCtx, span := obs.Tracer.Start(ctx, "getMenuSlug")
			defer span.End()

			menuID := c.PathParam("menuId")
			if menuID == "" {
				return MenuSlugResponse{}, fuego.HTTPError{
					Title:  "menuId is required",
					Detail: "menuId path parameter is required",
					Status: http.StatusBadRequest,
				}
			}

			userID, ok := sharedcontext.UserIDFromContext(spanCtx)
			if !ok || userID == "" {
				return MenuSlugResponse{}, fuego.HTTPError{
					Title:  "unauthorized",
					Detail: "user id not found in context",
					Status: http.StatusUnauthorized,
				}
			}

			hasMenu, err := userHasMenu(spanCtx, userID, menuID)
			if err != nil || !hasMenu {
				return MenuSlugResponse{}, fuego.HTTPError{
					Title:  "menu not found",
					Detail: "menu not found or you do not own it",
					Status: http.StatusNotFound,
				}
			}

			slug, err := getActiveSlug(spanCtx, menuID)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error getting active slug", "error", err, "menuID", menuID)
				return MenuSlugResponse{}, fuego.HTTPError{
					Title:  "error getting slug",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			if slug == "" {
				return MenuSlugResponse{}, fuego.HTTPError{
					Title:  "slug not found",
					Detail: "no active slug for this menu",
					Status: http.StatusNotFound,
				}
			}
			return MenuSlugResponse{Slug: slug}, nil
		},
		option.Summary("Get active slug for menu"),
		option.Description("Returns the active slug for the authenticated user's menu. Used by console for sharing."),
		option.Tags("menu", "slug"),
		option.Middleware(jwtAuthMiddleware),
	)
}
