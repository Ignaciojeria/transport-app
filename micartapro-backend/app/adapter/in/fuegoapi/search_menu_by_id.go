package fuegoapi

import (
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"net/http"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func init() {
	ioc.Register(searchMenuById)
}

func searchMenuById(
	s httpserver.Server,
	obs observability.Observability,
	getMenuById supabaserepo.GetMenuById,
) {
	fuego.Get(s.Manager, "/menu/{menuId}",
		func(c fuego.ContextNoBody) (events.MenuCreateRequest, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "searchMenuById")
			defer span.End()

			// Obtener el menu_id del parámetro de ruta
			menuID := c.PathParam("menuId")
			if menuID == "" {
				return events.MenuCreateRequest{}, fuego.HTTPError{
					Title:  "menuId is required",
					Detail: "menuId parameter is required",
					Status: http.StatusBadRequest,
				}
			}

			// Obtener el version_id opcional del query parameter
			versionID := c.QueryParam("version_id")

			// Obtener el menú desde Supabase usando el menu_id
			menu, err := getMenuById(spanCtx, menuID, versionID)
			if err != nil {
				if err == supabaserepo.ErrMenuNotFound {
					return events.MenuCreateRequest{}, fuego.HTTPError{
						Title:  "menu not found",
						Detail: "menu with the provided menu_id was not found",
						Status: http.StatusNotFound,
					}
				}
				obs.Logger.ErrorContext(spanCtx, "error getting menu from supabase by menu_id", "error", err)
				return events.MenuCreateRequest{}, fuego.HTTPError{
					Title:  "error getting menu",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			obs.Logger.InfoContext(spanCtx, "menu found successfully by menu_id", "menuID", menuID, "versionID", versionID)
			// Normalizar URLs de imágenes por si se guardaron mal en BD (httpshttps, https.storage, etc.)
			normalizeMenuImageURLs(&menu)
			return menu, nil
		},
		option.Summary("searchMenuById"),
		option.Tags("menu"),
		option.Path("menuId", "string", param.Required()))
}
