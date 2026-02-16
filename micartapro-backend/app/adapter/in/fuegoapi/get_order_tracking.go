package fuegoapi

import (
	"errors"

	"micartapro/app/adapter/out/supabaserepo"
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
		getOrderTracking,
		httpserver.New,
		observability.NewObservability,
		supabaserepo.NewGetOrderByTrackingID,
	)
}

func getOrderTracking(
	s httpserver.Server,
	obs observability.Observability,
	getOrderByTracking supabaserepo.GetOrderByTrackingID,
) {
	fuego.Get(s.Manager, "/api/tracking/{trackingId}",
		func(c fuego.ContextNoBody) (*supabaserepo.OrderByTrackingResult, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "getOrderTracking")
			defer span.End()

			trackingID := c.PathParam("trackingId")
			if trackingID == "" {
				return nil, fuego.HTTPError{
					Title:  "trackingId is required",
					Detail: "trackingId parameter is required",
					Status: http.StatusBadRequest,
				}
			}

			result, err := getOrderByTracking(spanCtx, trackingID)
			if err != nil {
				if errors.Is(err, supabaserepo.ErrOrderNotFound) {
					return nil, fuego.HTTPError{
						Title:  "order not found",
						Detail: "no order found for the provided tracking code",
						Status: http.StatusNotFound,
					}
				}
				obs.Logger.ErrorContext(spanCtx, "error getting order by tracking", "error", err, "trackingId", trackingID)
				return nil, fuego.HTTPError{
					Title:  "error getting order",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			return result, nil
		},
		option.Summary("getOrderTracking"),
		option.Tags("tracking"),
		option.Path("trackingId", "string", param.Required()))
}
