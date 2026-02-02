package fuegoapi

import (
	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

type UpdatePresentationStyleBody struct {
	PresentationStyle string `json:"presentationStyle"` // "HERO" | "MODERN"
}

func init() {
	ioc.Registry(
		updateMenuPresentationStyleHandler,
		httpserver.New,
		observability.NewObservability,
		supabaserepo.NewUpdateMenuPresentationStyle,
		apimiddleware.NewJWTAuthMiddleware,
	)
}

func updateMenuPresentationStyleHandler(
	s httpserver.Server,
	obs observability.Observability,
	updateStyle supabaserepo.UpdateMenuPresentationStyle,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware,
) {
	fuego.Patch(s.Manager, "/api/menus/{menuId}/presentation-style",
		func(c fuego.ContextWithBody[UpdatePresentationStyleBody]) (struct{}, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "updateMenuPresentationStyle")
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
			style := events.MenuPresentationStyle(body.PresentationStyle)
			if style != events.MenuStyleHero && style != events.MenuStyleModern {
				return struct{}{}, fuego.HTTPError{
					Title:  "invalid presentationStyle",
					Detail: "presentationStyle must be HERO or MODERN",
					Status: http.StatusBadRequest,
				}
			}

			if err := updateStyle(spanCtx, menuID, style); err != nil {
				if err == supabaserepo.ErrMenuNotFound {
					return struct{}{}, fuego.HTTPError{
						Title:  "menu not found",
						Detail: "menu with the provided menuId was not found",
						Status: http.StatusNotFound,
					}
				}
				obs.Logger.ErrorContext(spanCtx, "error updating menu presentation style", "error", err, "menuID", menuID)
				return struct{}{}, fuego.HTTPError{
					Title:  "error updating presentation style",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			return struct{}{}, nil
		},
		option.Summary("Update menu presentation style"),
		option.Description("Updates the presentation style (HERO or MODERN) of the menu's current version"),
		option.Tags("menu"),
		option.Middleware(jwtAuthMiddleware),
	)
}
