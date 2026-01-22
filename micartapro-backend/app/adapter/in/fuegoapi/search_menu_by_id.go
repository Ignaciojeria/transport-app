package fuegoapi

import (
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func init() {
	ioc.Registry(
		searchMenuById,
		httpserver.New,
		observability.NewObservability,
		supabaserepo.NewGetMenuBySlug,
	)
}

func searchMenuById(
	s httpserver.Server,
	obs observability.Observability,
	getMenuBySlug supabaserepo.GetMenuBySlug,
) {
	fuego.Get(s.Manager, "/menu/slug/{slug}",
		func(c fuego.ContextNoBody) (events.MenuCreateRequest, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "searchMenuById")
			defer span.End()

			// Obtener el slug del parámetro de ruta
			slug := c.PathParam("slug")
			if slug == "" {
				return events.MenuCreateRequest{}, fuego.HTTPError{
					Title:  "slug is required",
					Detail: "slug parameter is required",
					Status: http.StatusBadRequest,
				}
			}

			// Obtener el version_id opcional del query parameter
			versionID := c.QueryParam("version_id")

			// Obtener el menú desde Supabase usando el slug y opcionalmente el version_id
			menu, err := getMenuBySlug(spanCtx, slug, versionID)
			if err != nil {
				if err == supabaserepo.ErrMenuNotFound {
					return events.MenuCreateRequest{}, fuego.HTTPError{
						Title:  "menu not found",
						Detail: "menu with the provided slug was not found",
						Status: http.StatusNotFound,
					}
				}
				obs.Logger.ErrorContext(spanCtx, "error getting menu from supabase", "error", err)
				return events.MenuCreateRequest{}, fuego.HTTPError{
					Title:  "error getting menu",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			obs.Logger.InfoContext(spanCtx, "menu found successfully", "slug", slug, "menuID", menu.ID, "versionID", versionID)
			return menu, nil
		},
		option.Summary("searchMenuById"),
		option.Tags("menu"),
		option.Path("slug", "string", param.Required()))
}
