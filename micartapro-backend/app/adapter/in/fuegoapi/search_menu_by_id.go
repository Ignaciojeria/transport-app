package fuegoapi

import (
	"micartapro/app/adapter/out/storage"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
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
		supabaserepo.NewGetMenuSlugBySlug,
		storage.NewGetLatestMenuById,
	)
}

func searchMenuById(
	s httpserver.Server,
	obs observability.Observability,
	getMenuSlugBySlug supabaserepo.GetMenuSlugBySlug,
	getLatestMenuById storage.GetLatestMenuById,
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

			// Obtener user_id y menu_id desde el slug
			slugInfo, err := getMenuSlugBySlug(spanCtx, slug)
			if err != nil {
				if err == supabaserepo.ErrSlugNotFound {
					return events.MenuCreateRequest{}, fuego.HTTPError{
						Title:  "menu not found",
						Detail: "menu with the provided slug was not found",
						Status: http.StatusNotFound,
					}
				}
				obs.Logger.ErrorContext(spanCtx, "error getting menu slug", "error", err)
				return events.MenuCreateRequest{}, fuego.HTTPError{
					Title:  "error getting menu slug",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			// Agregar user_id al contexto para que GetLatestMenuById pueda usarlo
			spanCtx = sharedcontext.WithUserID(spanCtx, slugInfo.UserID)

			// Obtener el menú desde el storage
			menu, err := getLatestMenuById(spanCtx, slugInfo.MenuID)
			if err != nil {
				if err == storage.ErrMenuNotFound {
					return events.MenuCreateRequest{}, fuego.HTTPError{
						Title:  "menu not found",
						Detail: "menu was not found in storage",
						Status: http.StatusNotFound,
					}
				}
				obs.Logger.ErrorContext(spanCtx, "error getting menu from storage", "error", err)
				return events.MenuCreateRequest{}, fuego.HTTPError{
					Title:  "error getting menu from storage",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			obs.Logger.InfoContext(spanCtx, "menu found successfully", "slug", slug, "menuID", slugInfo.MenuID)
			return menu, nil
		},
		option.Summary("searchMenuById"),
		option.Tags("menu"),
		option.Path("slug", "string", param.Required()))
}
